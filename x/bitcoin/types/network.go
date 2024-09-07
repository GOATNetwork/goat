package types

import "github.com/btcsuite/btcd/chaincfg"

func (conf *ChainConfig) ToBtcdParam() *chaincfg.Params {
	// we only need the config for address codec
	return &chaincfg.Params{
		Name:             conf.NetworkName,
		PubKeyHashAddrID: uint8(conf.PubkeyHashAddrPrefix),
		ScriptHashAddrID: uint8(conf.ScriptHashAddrPrefix),
		Bech32HRPSegwit:  conf.Bech32Hrp,
	}
}
