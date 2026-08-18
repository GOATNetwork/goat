package keeper

import (
	"context"

	cmttypes "github.com/cometbft/cometbft/proto/tendermint/types"
)

// ConsensusParamsStore is an in memory stand in for x/consensus. The locking
// module writes to it when it widens the validator pub_key_types whitelist at
// a fork height, which is the only path there is: GOAT registers no x/gov and
// the ante handler admits no message that could carry a parameter update.
type ConsensusParamsStore struct {
	Params cmttypes.ConsensusParams
}

func NewConsensusParamsStore(pubKeyTypes ...string) *ConsensusParamsStore {
	return &ConsensusParamsStore{
		Params: cmttypes.ConsensusParams{
			Validator: &cmttypes.ValidatorParams{PubKeyTypes: pubKeyTypes},
		},
	}
}

func (c *ConsensusParamsStore) Get(context.Context) (cmttypes.ConsensusParams, error) {
	return c.Params, nil
}

func (c *ConsensusParamsStore) Set(_ context.Context, params cmttypes.ConsensusParams) error {
	c.Params = params
	return nil
}

// PubKeyTypes is a shorthand for the whitelist the store currently holds.
func (c *ConsensusParamsStore) PubKeyTypes() []string {
	if c.Params.Validator == nil {
		return nil
	}
	return c.Params.Validator.PubKeyTypes
}
