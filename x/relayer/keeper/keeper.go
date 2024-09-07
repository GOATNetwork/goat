package keeper

import (
	"context"
	"encoding/hex"
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
		Voters    collections.Map[string, types.Voter]
		Queue     collections.Item[types.VoterQueue]
		Pubkeys   collections.KeySet[[]byte]
		Randao    collections.Item[[]byte]

		accountKeeper types.AccountKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	addressCodec address.Codec,
	storeService store.KVStoreService,
	accountKeeper types.AccountKeeper,
	logger log.Logger,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:           cdc,
		AddrCodec:     addressCodec,
		storeService:  storeService,
		logger:        logger,
		accountKeeper: accountKeeper,

		Params:   collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Relayer:  collections.NewItem(sb, types.RelayerKey, "relayer", codec.CollValue[types.Relayer](cdc)),
		Sequence: collections.NewSequence(sb, types.SequenceKey, "sequence"),
		Voters:   collections.NewMap(sb, types.VotersKey, "voters", collections.StringKey, codec.CollValue[types.Voter](cdc)),
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
	proposer, err := k.Voters.Get(ctx, relayer.Proposer)
	if err != nil {
		return 0, err
	}
	pubkeys = append(pubkeys, proposer.VoteKey)

	for i := 0; i < len(voters); i++ {
		if !bmp.Contains(uint32(i)) {
			continue
		}

		voter, err := k.Voters.Get(ctx, voters[i])
		if err != nil {
			return 0, err
		}
		pubkeys = append(pubkeys, voter.VoteKey)
	}

	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	sigdoc := types.VoteSignDoc(req.MethodName(), sdkctx.ChainID(), relayer.Proposer, sequence, relayer.Epoch, req.VoteSigDoc())
	if !goatcrypto.AggregateVerify(pubkeys, sigdoc, req.GetVote().GetSignature()) {
		return 0, types.ErrInvalidProposalSignature.Wrapf("invalid signature")
	}

	for _, fn := range verifyFn {
		if err := fn(sigdoc); err != nil {
			return 0, err
		}
	}

	// As long as the proposer sends a valid tx, it should be considered that the proposer is accepted.
	if !relayer.ProposerAccepted {
		relayer.ProposerAccepted = true
		if err := k.Relayer.Set(ctx, relayer); err != nil {
			return 0, err
		}
	}

	return sequence, nil
}

func (k Keeper) VerifyNonProposal(ctx context.Context, req types.INonVoteMsg) (types.IRelayer, error) {
	relayer, err := k.Relayer.Get(ctx)
	if err != nil {
		return nil, err
	}

	if relayer.Proposer != req.GetProposer() {
		return nil, types.ErrNotProposer.Wrapf("not proposer")
	}

	// As long as the proposer sends a valid tx, it should be considered that the proposer is accepted.
	if !relayer.ProposerAccepted {
		relayer.ProposerAccepted = true
		if err := k.Relayer.Set(ctx, relayer); err != nil {
			return nil, err
		}
	}

	return &relayer, nil
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

// ElectProposer elects new proposer and increase epoch number at regular intervals
func (k Keeper) ElectProposer(ctx context.Context) error {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)

	relayer, err := k.Relayer.Get(ctx)
	if err != nil {
		return err
	}

	param, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	if duration := sdkctx.BlockTime().Sub(relayer.LastElected); duration < param.ElectingPeriod &&
		(relayer.ProposerAccepted || param.AcceptProposerTimeout == 0 || duration < param.AcceptProposerTimeout) {
		return nil
	}

	queue, err := k.Queue.Get(ctx)
	if err != nil {
		return err
	}

	onBoarding, offBoarding := len(queue.OnBoarding) > 0, len(queue.OffBoarding) > 0
	if onBoarding {
		for _, v := range queue.OnBoarding {
			voter, err := k.Voters.Get(ctx, v)
			if err != nil {
				return err
			}
			voter.Status = types.VOTER_STATUS_ACTIVATED
			if err := k.Voters.Set(ctx, v, voter); err != nil {
				return err
			}
		}
		relayer.Voters = append(relayer.Voters, queue.OnBoarding...)
	}

	var isProposerRemvoed bool
	if offBoarding {
		set := make(map[string]bool, len(queue.OffBoarding))
		for _, v := range queue.OffBoarding {
			if err := k.Voters.Remove(ctx, v); err != nil {
				return err
			}
			set[v] = true
		}

		isProposerRemvoed = set[relayer.Proposer]
		newVoters := slices.DeleteFunc(relayer.Voters, func(e string) bool {
			return set[e]
		})

		if isProposerRemvoed {
			if len(newVoters) == 0 { // it should never happen
				k.Logger().Error("delete too many voters in ElectProposer")
				return nil
			}
			// use the first voter as the new proposer
			relayer.Proposer = newVoters[0]
			relayer.Voters = newVoters[1:]
		} else {
			relayer.Voters = newVoters[:]
		}
	}

	// epoch number is constantly increasing even if we don't have a new election
	relayer.Epoch++
	relayer.LastElected = sdkctx.BlockTime()

	var events = sdktypes.Events{types.NewEpochEvent(relayer.Epoch)}
	if offBoarding || onBoarding {
		events = append(events, types.VoterChangedEvent(relayer.Epoch, queue.OnBoarding, queue.OffBoarding)...)

		queue.OnBoarding = queue.OnBoarding[:0]
		queue.OffBoarding = queue.OffBoarding[:0]
		if err := k.Queue.Set(ctx, queue); err != nil {
			return err
		}

		// if the proposer is removed, we don't make a election, just use the next voter as the new proposer
		if isProposerRemvoed {
			relayer.ProposerAccepted = false

			k.Logger().Debug("New proposer", "height", sdkctx.BlockHeight(), "epoch", relayer.Epoch, "proposer", relayer.Proposer)
			if err := k.Relayer.Set(ctx, relayer); err != nil {
				return err
			}

			sdkctx.EventManager().EmitEvents(
				append(events, types.ElectedProposerEvent(relayer.Proposer, relayer.Epoch)),
			)
			return nil
		}
	}

	voterLen := len(relayer.Voters)
	// no voter no election
	if voterLen == 0 {
		relayer.ProposerAccepted = true
		if err := k.Relayer.Set(ctx, relayer); err != nil {
			return err
		}
		sdkctx.EventManager().EmitEvents(events)
		return nil
	}

	// only get hash when we have 2 voters at least
	if voterLen > 1 {
		randao, err := k.Randao.Get(ctx)
		if err != nil {
			return err
		}
		// hash with the current epoch to ensure always have a new randao value
		rand := new(big.Int).SetBytes(goatcrypto.SHA256Sum(randao, goatcrypto.Uint64LE(relayer.Epoch)))
		proposerIndex := rand.Mod(rand, big.NewInt(int64(voterLen))).Int64()
		relayer.Proposer, relayer.Voters[proposerIndex] = relayer.Voters[proposerIndex], relayer.Proposer
	} else {
		relayer.Proposer, relayer.Voters[0] = relayer.Voters[0], relayer.Proposer
	}

	k.Logger().Debug("New proposer", "height", sdkctx.BlockHeight(), "epoch", relayer.Epoch, "proposer", relayer.Proposer)
	relayer.ProposerAccepted = false
	if err := k.Relayer.Set(ctx, relayer); err != nil {
		return err
	}

	sdkctx.EventManager().EmitEvents(
		append(events, types.ElectedProposerEvent(relayer.Proposer, relayer.Epoch)),
	)
	return nil
}

func (k Keeper) GetCurrentProposer(ctx context.Context) (sdktypes.AccAddress, error) {
	relayer, err := k.Relayer.Get(ctx)
	if err != nil {
		return nil, err
	}
	return k.AddrCodec.StringToBytes(relayer.Proposer)
}
