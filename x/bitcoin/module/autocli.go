package bitcoin

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	modulev1 "github.com/goatnetwork/goat/api/goat/bitcoin/v1"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "NewBlockHashes",
					Skip:      true,
				},
				{
					RpcMethod: "NewDeposits",
					Skip:      true,
				},
				{
					RpcMethod: "NewPubkey",
					Skip:      true,
				},
				{
					RpcMethod: "ProcessWithdrawal",
					Skip:      true,
				},
				{
					RpcMethod: "ReplaceWithdrawal",
					Skip:      true,
				},
				{
					RpcMethod: "FinalizeWithdrawal",
					Skip:      true,
				},
				{
					RpcMethod: "ApproveCancellation",
					Skip:      true,
				},
				{
					RpcMethod: "FinalizeWithdrawal",
					Skip:      true,
				},
			},
		},
	}
}
