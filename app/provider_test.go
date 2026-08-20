package app

import (
	"os"
	"path/filepath"
	"testing"

	"cosmossdk.io/log/v2"
	cmtcrypto "github.com/cometbft/cometbft/crypto"
	cmmldsa65 "github.com/cometbft/cometbft/crypto/mldsa65"
	cmsecp256k1 "github.com/cometbft/cometbft/crypto/secp256k1"
	cmtjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/privval"
	bitcintypes "github.com/goatnetwork/goat/x/bitcoin/types"
	"github.com/stretchr/testify/require"
)

type testAppOptions map[string]any

func (o testAppOptions) Get(key string) any { return o[key] }

// writeKey stores priv in the privval file format, which is what both the
// consensus key and the split transaction signing key use.
func writeKey(t *testing.T, path string, priv cmtcrypto.PrivKey) {
	t.Helper()
	key := privval.FilePVKey{
		Address: priv.PubKey().Address(),
		PubKey:  priv.PubKey(),
		PrivKey: priv,
	}
	data, err := cmtjson.MarshalIndent(&key, "", "  ")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(path, data, 0o600))
}

func mldsaKey(t *testing.T) cmtcrypto.PrivKey {
	t.Helper()
	priv, err := cmmldsa65.GenPrivKey()
	require.NoError(t, err)
	return priv
}

// setup lays out a node home and returns the options plus the two key paths.
func setup(t *testing.T) (testAppOptions, string, string) {
	t.Helper()
	home := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(home, "config"), 0o700))
	opts := testAppOptions{
		"home":                    home,
		"priv_validator_key_file": filepath.Join("config", "priv_validator_key.json"),
	}
	return opts,
		filepath.Join(home, "config", "priv_validator_key.json"),
		filepath.Join(home, "config", TxSigningKeyFile)
}

// A validator that has never rotated has no split key file, and the privval
// key stays the account key. This is every node on the network today.
func TestFallsBackToPrivval(t *testing.T) {
	opts, privvalPath, _ := setup(t)
	consensus := cmsecp256k1.GenPrivKey()
	writeKey(t, privvalPath, consensus)

	got := ProvideValidatorPrvKey(log.NewNopLogger(), opts)
	require.Equal(t, consensus.Bytes(), got.Bytes())
}

// Once split, the account key comes from the dedicated file and the privval
// key is free to be post quantum. Without this the node cannot start at all.
func TestSplitKeyWinsOverMlDsaPrivval(t *testing.T) {
	opts, privvalPath, txKeyPath := setup(t)
	account := cmsecp256k1.GenPrivKey()
	writeKey(t, privvalPath, mldsaKey(t))
	writeKey(t, txKeyPath, account)

	got := ProvideValidatorPrvKey(log.NewNopLogger(), opts)
	require.Equal(t, account.Bytes(), got.Bytes())
	require.Equal(t, bitcintypes.Secp256K1Name, got.Type())
}

// The split file is preferred even when the privval key would also have been
// accepted, so that a rotation to another secp256k1 key does not silently
// replace the account key.
func TestSplitKeyWinsOverSecp256k1Privval(t *testing.T) {
	opts, privvalPath, txKeyPath := setup(t)
	account := cmsecp256k1.GenPrivKey()
	writeKey(t, privvalPath, cmsecp256k1.GenPrivKey())
	writeKey(t, txKeyPath, account)

	got := ProvideValidatorPrvKey(log.NewNopLogger(), opts)
	require.Equal(t, account.Bytes(), got.Bytes())
}

// Rotating the privval file without splitting first is the mistake this
// panic exists to name, so it has to say where the original key belongs.
func TestMlDsaPrivvalWithoutSplitPanics(t *testing.T) {
	opts, privvalPath, txKeyPath := setup(t)
	writeKey(t, privvalPath, mldsaKey(t))

	require.PanicsWithValue(t,
		privvalPath+" is not an secp256k1 key. A validator that has rotated its"+
			" consensus key has to keep the key its account was created with in "+txKeyPath,
		func() { ProvideValidatorPrvKey(log.NewNopLogger(), opts) })
}

// A split file of the wrong type is a different mistake and names itself.
func TestMlDsaSplitKeyPanics(t *testing.T) {
	opts, privvalPath, txKeyPath := setup(t)
	writeKey(t, privvalPath, cmsecp256k1.GenPrivKey())
	writeKey(t, txKeyPath, mldsaKey(t))

	require.PanicsWithValue(t, txKeyPath+" is not an secp256k1 key",
		func() { ProvideValidatorPrvKey(log.NewNopLogger(), opts) })
}

func TestMissingPrivvalPanics(t *testing.T) {
	opts, _, _ := setup(t)
	require.Panics(t, func() { ProvideValidatorPrvKey(log.NewNopLogger(), opts) })
}

func TestKeyFileWithoutPrivateKeyPanics(t *testing.T) {
	opts, privvalPath, _ := setup(t)
	require.NoError(t, os.WriteFile(privvalPath, []byte(`{}`), 0o600))
	require.Panics(t, func() { ProvideValidatorPrvKey(log.NewNopLogger(), opts) })
}
