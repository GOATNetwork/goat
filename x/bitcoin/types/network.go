package types

import "github.com/btcsuite/btcd/chaincfg"

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
