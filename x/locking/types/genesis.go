package types

import (
	"errors"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params:  DefaultParams(),
		Slashed: sdk.NewCoins(),
		RewardPool: RewardPool{
			Goat:   math.ZeroInt(),
			Gas:    math.ZeroInt(),
			Remain: math.ZeroInt(),
		},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}
	if err := gs.Slashed.Validate(); err != nil {
		return err
	}

	if gs.RewardPool.Goat.IsNegative() {
		return errors.New("reward pool goat amount cannot be negative")
	}
	if gs.RewardPool.Gas.IsNegative() {
		return errors.New("reward pool gas amount cannot be negative")
	}
	if gs.RewardPool.Remain.IsNegative() {
		return errors.New("reward pool remain amount cannot be negative")
	}

	if len(gs.Tokens) == 0 {
		return errors.New("no tokens in genesis")
	}
	for i, val := range gs.Tokens {
		if err := sdk.ValidateDenom(val.Denom); err != nil {
			return fmt.Errorf("invalid token denom: index %d", i)
		}
		if val.Token.Threshold.IsNegative() {
			return fmt.Errorf("token threshold cannot be negative: index %d", i)
		}
	}

	if len(gs.Validators) == 0 {
		return errors.New("no validators in genesis")
	}
	for i, val := range gs.Validators {
		if len(val.Pubkey) != 33 {
			return fmt.Errorf("invalid pubkey length: index %d", i)
		}
		if first := val.Pubkey[0]; first != 2 && first != 3 {
			return fmt.Errorf("invalid pubkey prefix: index %d", i)
		}

		if err := val.Locking.Validate(); err != nil {
			return err
		}
		if val.Reward.IsNegative() {
			return fmt.Errorf("validator reward cannot be negative: index %d", i)
		}
		if val.GasReward.IsNegative() {
			return fmt.Errorf("validator gas reward cannot be negative: index %d", i)
		}
		if val.SigningInfo.Missed < 0 || val.SigningInfo.Offset < 0 {
			return fmt.Errorf("invalid signing info values: index %d", i)
		}
	}

	for i, reward := range gs.EthTxQueue.Rewards {
		if len(reward.Recipient) != 20 {
			return fmt.Errorf("invalid eth tx queue reward recipient length: index %d", i)
		}
		if reward.Goat.IsNegative() || reward.Gas.IsNegative() {
			return fmt.Errorf("eth tx queue reward amount cannot be negative: index %d", i)
		}
	}

	for i, tx := range gs.EthTxQueue.Unlocks {
		if len(tx.Recipient) != 20 {
			return fmt.Errorf("invalid eth tx queue unlock recipient length: index %d", i)
		}
		if len(tx.Token) != 20 {
			return fmt.Errorf("invalid eth tx queue unlock token length: index %d", i)
		}
		if tx.Amount.IsNegative() {
			return fmt.Errorf("eth tx queue unlock amount cannot be negative: index %d", i)
		}
	}

	for i, val := range gs.UnlockQueue {
		for j, lock := range val.Unlocks {
			if len(lock.Token) != 20 {
				return fmt.Errorf("invalid unlock queue lock token length: val index %d, lock index %d", i, j)
			}
			if len(lock.Recipient) != 20 {
				return fmt.Errorf("invalid unlock queue lock recipient length: val index %d, lock index %d", i, j)
			}
			if lock.Amount.IsNegative() {
				return fmt.Errorf("unlock queue lock amount cannot be negative: val index %d, lock index %d", i, j)
			}
		}
	}

	return nil
}
