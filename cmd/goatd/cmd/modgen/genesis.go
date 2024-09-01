package modgen

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/gogoproto/proto"
)

func UpdateGensis[T proto.Message](path, module string, state T, codec codec.JSONCodec, updateFn func(state T) error) error {
	if updateFn == nil {
		return errors.New("updateFn is nil")
	}

	origin, err := types.AppGenesisFromFile(path)
	if err != nil {
		return err
	}

	var appGenState map[string]json.RawMessage
	if err := json.Unmarshal(origin.AppState, &appGenState); err != nil {
		return err
	}

	// extract the module gensis
	rawState, ok := appGenState[module]
	if !ok {
		return fmt.Errorf("%s doesn't exist in the genesis file", module)
	}

	if err := codec.UnmarshalJSON(rawState, state); err != nil {
		return err
	}

	// update the module genesis
	if err := updateFn(state); err != nil {
		return err
	}

	// encode the module genesis
	rawState, err = codec.MarshalJSON(state)
	if err != nil {
		return err
	}
	appGenState[module] = rawState

	// update the app gensis field
	newAppState, err := json.Marshal(appGenState)
	if err != nil {
		return err
	}
	origin.AppState = newAppState

	// wirte back to the genesis file
	return genutil.ExportGenesisFile(origin, path)
}
