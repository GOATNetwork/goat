package modgen

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
)

const (
	FlagRegtest = "regtest"
)

func DecodeHexOrBase64String(str string) ([]byte, error) {
	if strings.HasPrefix(str, "0x") {
		return hex.DecodeString(str[2:])
	}
	pubkeyRaw, err := hex.DecodeString(str)
	if err != nil {
		pubkeyRaw, err = base64.StdEncoding.DecodeString(str)
		if err != nil {
			return nil, fmt.Errorf("pubkey %s doesn't use base64 or hex encoding", str)
		}
	}
	return pubkeyRaw, nil
}

func GetCompressedK256P1Pubkey(pubkeyRaw []byte) ([]byte, error) {
	switch len(pubkeyRaw) {
	case 33:
		if pubkeyRaw[0] != 2 && pubkeyRaw[0] != 3 {
			return nil, errors.New("invalid compressed secp256k1 pubkey prefix")
		}
	case 64, 65:
		if len(pubkeyRaw) == 65 {
			if pubkeyRaw[0] != 4 {
				return nil, errors.New("invalid uncompressed secp256k1 pubkey prefix")
			}
			pubkeyRaw = pubkeyRaw[1:]
		}
		pubkeyRaw = goatcrypto.CompressP256k1Pubkey([64]byte(pubkeyRaw))
	}
	return pubkeyRaw, nil
}

func GetEthGenesisHeaderByFile(configPath string) (*ethtypes.Header, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	type Config struct {
		Consensus struct {
			Goat *ethtypes.Header
		}
	}

	genesis := new(Config)
	if err := json.NewDecoder(file).Decode(genesis); err != nil {
		return nil, err
	}

	header := genesis.Consensus.Goat
	if header == nil {
		return nil, errors.New("genesis header is empty")
	}

	if header.BaseFee == nil || header.WithdrawalsHash == nil {
		return nil, errors.New("shanghai upgrade should be activated")
	}

	if *header.WithdrawalsHash != ethtypes.EmptyWithdrawalsHash {
		return nil, errors.New("no withdrawals required")
	}

	if header.GasUsed != 0 || header.TxHash != ethtypes.EmptyTxsHash {
		return nil, errors.New("no txs required")
	}

	if header.BlobGasUsed == nil || header.ExcessBlobGas == nil || header.ParentBeaconRoot == nil {
		return nil, errors.New("cancun upgrade should be activated")
	}

	if *header.BlobGasUsed != 0 {
		return nil, errors.New("required no blob txes")
	}

	if header.RequestsHash == nil {
		return nil, errors.New("no requests provided")
	}

	return header, nil
}

func getValidatorSignMsg(chainID uint64, owner, validator []byte) []byte {
	data := new(big.Int).SetUint64(chainID).FillBytes(make([]byte, 32))
	return ethcrypto.Keccak256(data, validator, owner)
}

func PrintJSON(info any) error {
	out, err := json.MarshalIndent(info, "", " ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(os.Stdout, "%s\n", out)
	return err
}

func PrintStderr(a ...any) {
	fmt.Fprintln(os.Stderr, a...)
}
