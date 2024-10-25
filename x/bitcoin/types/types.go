package types

const (
	NewPubkeyMethodSigName         = "Bitcoin/NewPubkey"
	NewBlocksMethodSigName         = "Bitcoin/NewBlocks"
	ProcessWithdrawalMethodSigName = "Bitcoin/ProcessWithdrawal"
	ReplaceWithdrawalMethodSigName = "Bitcoin/ReplaceWithdrawal"
	NewConsolidationMethodSigName  = "Bitcoin/NewConsolidation"
)

const (
	DepositMagicLen    = 4
	DustTxoutAmount    = 546
	RawBtcHeaderSize   = 80
	P2TRScriptSize     = 34
	P2WSHScriptSize    = 34
	P2WPKHScriptSize   = 22
	MinBtcTxSize       = 4 + 1 + 32 + 4 + 1 + 0 + 4 + 1 + 8 + 1 + 22 + 4
	DepositV1TxoutSize = 26
	// 4 version
	// 1 input length
	// 41 = 32 prevTxid + 4 prevTxOut + 1 sigScriptLength + 0 sigScript(witness) + 4 sequence
	// 1 output length
	// 8 value + 1 pkScriptLength + 34 p2wsh/p2tr
	// || 8 value + 1 pkScriptLength + 22 p2wph +  8 value + 1 pkScriptLength + 26 data OP_RETURN
	// 4 lockTime
	MinDepositTxSize    = 4 + 1 + 32 + 4 + 1 + 0 + 4 + 1 + 8 + 1 + 34 + 4
	MaxAllowedBtcTxSize = 32 * 1024
)
