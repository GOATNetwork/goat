package types

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		new(MsgNewBlockHashes),
		new(MsgNewDeposits),
		new(MsgNewPubkey),
		new(MsgProcessWithdrawal),
		new(MsgProcessWithdrawalV2),
		new(MsgReplaceWithdrawal),
		new(MsgReplaceWithdrawalV2),
		new(MsgFinalizeWithdrawal),
		new(MsgApproveCancellation),
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
