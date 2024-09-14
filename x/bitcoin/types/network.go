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

var BitcoinNetworks = map[string]*chaincfg.Params{
	chaincfg.MainNetParams.Name:       &chaincfg.MainNetParams,
	chaincfg.TestNet3Params.Name:      &chaincfg.TestNet3Params,
	chaincfg.SigNetParams.Name:        &chaincfg.SigNetParams,
	chaincfg.RegressionNetParams.Name: &chaincfg.RegressionNetParams,
}

var DepositMagicPreifxs = map[string]string{
	chaincfg.MainNetParams.Name:       "GTV2",
	chaincfg.TestNet3Params.Name:      "GTV1",
	chaincfg.SigNetParams.Name:        "GTV1",
	chaincfg.RegressionNetParams.Name: "GTT0",
}
