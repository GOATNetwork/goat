package keeper

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/kelindar/bitmap"

	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
	"github.com/goatnetwork/goat/pkg/utils"
	"github.com/goatnetwork/goat/x/relayer/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		addressCodec address.Codec
		storeService store.KVStoreService
		logger       log.Logger

		// the address capable of executing a MsgUpdateParams message.
		// Typically, this should be the x/gov module account.
		authority string

		schema      collections.Schema
		Params      collections.Item[types.Params]
		Relayer     collections.Item[types.Relayer]
		Epoch       collections.Sequence
		ProposalSeq collections.Sequence
		Voters      collections.Map[string, types.Voter]
		Pubkeys     collections.KeySet[[]byte]
		Randao      collections.Item[[]byte]
		// this line is used by starport scaffolding # collection/type

	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	addressCodec address.Codec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,

) Keeper {
	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:          cdc,
		addressCodec: addressCodec,
		storeService: storeService,
		authority:    authority,
		logger:       logger,

		Params:      collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Relayer:     collections.NewItem(sb, types.RelayerKey, "relayer", codec.CollValue[types.Relayer](cdc)),
		Epoch:       collections.NewSequence(sb, types.EpochKey, "epoch"),
		ProposalSeq: collections.NewSequence(sb, types.ProposalKey, "proposal"),
		Voters:      collections.NewMap(sb, types.VotersKey, "voters", collections.StringKey, codec.CollValue[types.Voter](cdc)),
		Pubkeys:     collections.NewKeySet(sb, types.PubkeysKey, "pubkeys", collections.BytesKey),
		Randao:      collections.NewItem(sb, types.RandDAOKey, "randao", collections.BytesValue),
		// this line is used by starport scaffolding # collection/instantiate
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.schema = schema

	return k
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) VerifyProposal(ctx context.Context, req types.IVoteMsg, verifyFn ...func(sigdoc []byte) error) (uint64, error) {
	relayer, err := k.Relayer.Get(ctx)
	if err != nil {
		return 0, err
	}

	if relayer.Proposer != req.GetProposer() {
		return 0, types.ErrNotProposer.Wrapf("not proposer")
	}

	requireVoters := relayer.GetVoters()
	requestVoters := req.GetVote()
	if requestVoters == nil {
		return 0, types.ErrInvalidProposalSignature.Wrap("no vote info")
	}

	sequence, err := k.ProposalSeq.Peek(ctx)
	if err != nil {
		return 0, err
	}

	if requestVoters.GetSequence() != sequence {
		return 0, types.ErrInvalidProposalSignature.Wrap("incorrect seqeuence")
	}

	voterBitmap := bitmap.FromBytes(requestVoters.GetVoters())

	voterLen := voterBitmap.Count()
	if voterLen+1 < int(relayer.Threshold) || voterLen > len(requireVoters) {
		return 0, types.ErrInvalidProposalSignature.Wrapf("malformed signature length")
	}

	pubkeys := make([][]byte, 0, voterLen+1)
	proposer, err := k.Voters.Get(ctx, relayer.Proposer)
	if err != nil {
		return 0, err
	}
	pubkeys = append(pubkeys, proposer.VoteKey)

	for i := 0; i < len(requireVoters); i++ {
		if !voterBitmap.Contains(uint32(i)) {
			continue
		}

		voter, err := k.Voters.Get(ctx, relayer.GetVoters()[i])
		if err != nil {
			return 0, err
		}
		pubkeys = append(pubkeys, voter.VoteKey)
	}

	chainId := sdktypes.UnwrapSDKContext(ctx).HeaderInfo().ChainID

	sigdoc := types.VoteSignDoc(req.MethodName(), chainId, relayer.Proposer, sequence, req.VoteSigDoc())
	if !goatcrypto.AggregateVerify(pubkeys, sigdoc, requestVoters.GetSignature()) {
		return 0, types.ErrInvalidProposalSignature.Wrapf("invalid signature")
	}

	if len(verifyFn) == 1 && verifyFn[0] != nil {
		return sequence, verifyFn[0](sigdoc)
	}

	return sequence, nil
}

func (k Keeper) UpdateRandao(ctx context.Context, req types.IVoteMsg) error {
	randao, err := k.Randao.Get(ctx)
	if err != nil {
		return err
	}

	newRandao := utils.SHA256Sum(randao, req.GetVote().Signature)
	if err := k.Randao.Set(ctx, newRandao); err != nil {
		return err
	}
	k.Logger().Debug("Randao updated", "previous", hex.EncodeToString(randao), "current", hex.EncodeToString(newRandao))
	return nil
}

func (k Keeper) HasPubkey(ctx context.Context, raw []byte) (bool, error) {
	return k.Pubkeys.Has(ctx, raw)
}

func (k Keeper) AddNewKey(ctx context.Context, raw []byte) error {
	return k.Pubkeys.Set(ctx, raw)
}

func (k Keeper) SetProposalSeq(ctx context.Context, seq uint64) error {
	return k.ProposalSeq.Set(ctx, seq)
}

func (k Keeper) ElecteProposer(ctx context.Context) error {
	relayer, err := k.Relayer.Get(ctx)
	if err != nil {
		return err
	}

	sdkctx := sdktypes.UnwrapSDKContext(ctx)

	header := sdkctx.HeaderInfo()
	if vlen := int64(len(relayer.Voters)); header.Time.Sub(relayer.LastElected) > 10*time.Minute && vlen > 0 {
		epoch, err := k.Epoch.Peek(ctx)
		if err != nil {
			return err
		}

		randao, err := k.Randao.Get(ctx)
		if err != nil {
			return err
		}

		epoch++
		epochRaw := make([]byte, 8)
		binary.LittleEndian.PutUint64(epochRaw, epoch)

		// hash with the current epoch to ensure always have a new randao value
		curRand := new(big.Int).SetBytes(utils.SHA256Sum(randao, epochRaw))
		proposerIndex := curRand.Mod(curRand, big.NewInt(vlen)).Int64()

		relayer.Proposer, relayer.Voters[proposerIndex] = relayer.Voters[proposerIndex], relayer.Proposer

		relayer.Version++
		relayer.LastElected = header.Time

		k.Logger().Debug("New proposer", "height", header.Height, "proposer", relayer.Proposer)
		if err := k.Relayer.Set(ctx, relayer); err != nil {
			return err
		}
		if err := k.Epoch.Set(ctx, epoch); err != nil {
			return err
		}
		sdkctx.EventManager().EmitEvent(types.NewProposer(relayer.Proposer))
	}

	return nil
}
