package keeper

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"slices"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/kelindar/bitmap"

	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
	"github.com/goatnetwork/goat/x/relayer/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger
		schema       collections.Schema

		AddrCodec address.Codec
		Params    collections.Item[types.Params]
		Relayer   collections.Item[types.Relayer]
		Sequence  collections.Sequence
		Voters    collections.Map[sdktypes.AccAddress, types.Voter]
		Queue     collections.Item[types.VoterQueue]
		Pubkeys   collections.KeySet[[]byte]
		Randao    collections.Item[[]byte]
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	addressCodec address.Codec,
	storeService store.KVStoreService,
	logger log.Logger,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:          cdc,
		AddrCodec:    addressCodec,
		storeService: storeService,
		logger:       logger,

		Params:   collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Relayer:  collections.NewItem(sb, types.RelayerKey, "relayer", codec.CollValue[types.Relayer](cdc)),
		Sequence: collections.NewSequence(sb, types.SequenceKey, "sequence"),
		Voters:   collections.NewMap(sb, types.VotersKey, "voters", sdktypes.AccAddressKey, codec.CollValue[types.Voter](cdc)),
		Pubkeys:  collections.NewKeySet(sb, types.PubkeysKey, "pubkeys", collections.BytesKey),
		Queue:    collections.NewItem(sb, types.QueueKey, "queue", codec.CollValue[types.VoterQueue](cdc)),
		Randao:   collections.NewItem(sb, types.RandDAOKey, "randao", collections.BytesValue),
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

	voters := relayer.GetVoters()
	sequence, err := k.Sequence.Peek(ctx)
	if err != nil {
		return 0, err
	}

	if req.GetVote().GetSequence() != sequence {
		return 0, types.ErrInvalidProposalSignature.Wrap("incorrect seqeuence")
	}

	if req.GetVote().GetEpoch() != relayer.Epoch {
		return 0, types.ErrInvalidProposalSignature.Wrap("incorrect epoch")
	}

	threshold := (1 + len(voters)) * 2 / 3
	if threshold == 0 {
		threshold = 1
	}

	bmp := bitmap.FromBytes(req.GetVote().GetVoters())
	bmpLen := bmp.Count()
	if bmpLen+1 < threshold || bmpLen > len(voters) {
		return 0, types.ErrInvalidProposalSignature.Wrapf("malformed signature length")
	}

	pubkeys := make([][]byte, 0, bmpLen+1)
	p, err := k.AddrCodec.StringToBytes(relayer.Proposer)
	if err != nil {
		return 0, err
	}

	proposer, err := k.Voters.Get(ctx, p)
	if err != nil {
		return 0, err
	}
	pubkeys = append(pubkeys, proposer.VoteKey)

	for i := 0; i < len(voters); i++ {
		if !bmp.Contains(uint32(i)) {
			continue
		}

		v, err := k.AddrCodec.StringToBytes(relayer.GetVoters()[i])
		if err != nil {
			return 0, err
		}

		voter, err := k.Voters.Get(ctx, v)
		if err != nil {
			return 0, err
		}
		pubkeys = append(pubkeys, voter.VoteKey)
	}

	chainId := sdktypes.UnwrapSDKContext(ctx).HeaderInfo().ChainID

	sigdoc := types.VoteSignDoc(req.MethodName(), chainId, relayer.Proposer, sequence, relayer.Epoch, req.VoteSigDoc())
	if !goatcrypto.AggregateVerify(pubkeys, sigdoc, req.GetVote().GetSignature()) {
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

	newRandao := goatcrypto.SHA256Sum(randao, req.GetVote().Signature)
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
	return k.Sequence.Set(ctx, seq)
}

func (k Keeper) ElecteProposer(ctx context.Context) error {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)

	relayer, err := k.Relayer.Get(ctx)
	if err != nil {
		return err
	}

	param, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	queue, err := k.Queue.Get(ctx)
	if err != nil {
		return err
	}

	onBoarding, offBoarding := len(queue.OnBoarding) > 0, len(queue.OffBoarding) > 0
	if onBoarding {
		for _, v := range queue.OnBoarding {
			addr, err := k.AddrCodec.StringToBytes(v)
			if err != nil {
				return err
			}
			voter, err := k.Voters.Get(ctx, addr)
			if err != nil {
				return err
			}
			voter.Status = types.VOTER_STATUS_ACTIVATED
		}
		relayer.Voters = append(relayer.Voters, queue.OnBoarding...)
	}

	if offBoarding {
		set := make(map[string]bool)
		for _, v := range queue.OffBoarding {
			addr, err := k.AddrCodec.StringToBytes(v)
			if err != nil {
				return err
			}
			if err := k.Voters.Remove(ctx, addr); err != nil {
				return err
			}
			set[v] = true
		}

		voters := slices.DeleteFunc(append([]string{relayer.Proposer}, relayer.Voters...), func(e string) bool {
			return set[e]
		})

		if len(voters) == 0 { // it should never happen
			return errors.New("Too few voters avaliable")
		}

		relayer.Proposer = voters[0]
		relayer.Voters = voters[1:]
	}

	if offBoarding || onBoarding {
		queue.OnBoarding = queue.OnBoarding[:0]
		queue.OffBoarding = queue.OffBoarding[:0]
		if err := k.Queue.Set(ctx, queue); err != nil {
			return err
		}
	}

	blockTime := sdkctx.BlockTime()
	vlen, duration := int64(len(relayer.Voters)), blockTime.Sub(relayer.LastElected)
	acceptTimeout := param.AcceptProposerTimeout > 0 && duration > param.AcceptProposerTimeout
	if vlen > 0 && (duration > param.ElectingPeriod || acceptTimeout) {
		randao, err := k.Randao.Get(ctx)
		if err != nil {
			return err
		}

		relayer.Epoch++

		// hash with the current epoch to ensure always have a new randao value
		curRand := new(big.Int).SetBytes(goatcrypto.SHA256Sum(randao, goatcrypto.Uint64LE(relayer.Epoch)))
		proposerIndex := curRand.Mod(curRand, big.NewInt(vlen)).Int64()

		relayer.Proposer, relayer.Voters[proposerIndex] = relayer.Voters[proposerIndex], relayer.Proposer
		relayer.LastElected = blockTime
		relayer.ProposerAccepted = false

		k.Logger().Debug("New proposer", "height", sdkctx.BlockHeight(), "proposer", relayer.Proposer)
		if err := k.Relayer.Set(ctx, relayer); err != nil {
			return err
		}
		sdkctx.EventManager().EmitEvent(types.ElectedProposerEvent(relayer.Proposer, relayer.Epoch))
	}

	return nil
}

func (k Keeper) GetCurrentProposer(ctx context.Context) (sdktypes.AccAddress, error) {
	relayer, err := k.Relayer.Get(ctx)
	if err != nil {
		return nil, err
	}
	return k.AddrCodec.StringToBytes(relayer.Proposer)
}
