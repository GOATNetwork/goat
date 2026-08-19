package cmd

import (
	"bytes"
	"encoding/hex"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/cometbft/cometbft/privval"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/goatnetwork/goat/x/locking/types"
)

const (
	testChainID   = "goat-mainnet"
	testValidator = "0x1234567890123456789012345678901234567890"
)

var fieldRe = regexp.MustCompile(`(?m)^(\w+)\s+(?:0x)?(\S+)`)

func runRotateProof(t *testing.T, args ...string) map[string]string {
	t.Helper()
	cmd := RotateProofCommand()
	out := new(bytes.Buffer)
	cmd.SetOut(out)
	cmd.SetErr(out)
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v (%s)", err, out)
	}

	fields := map[string]string{}
	for _, m := range fieldRe.FindAllStringSubmatch(out.String(), -1) {
		fields[m[1]] = m[2]
	}
	return fields
}

func decode(t *testing.T, s string) []byte {
	t.Helper()
	raw, err := hex.DecodeString(s)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

// The whole point of the command: whatever it prints has to be something
// VerifyRotationProof accepts, or the rotation is soft-failed on chain with
// only a log line to show for it.
func TestRotateProofProducesAVerifiableProof(t *testing.T) {
	out := filepath.Join(t.TempDir(), "new_priv_validator_key.json")
	fields := runRotateProof(t,
		"--chain-id", testChainID, "--validator", testValidator, "--out", out)

	pubkey, proof := decode(t, fields["pubkey"]), decode(t, fields["proof"])
	if len(pubkey) != 1952 || len(proof) != 3309 {
		t.Fatalf("ml_dsa_65 sizes: pubkey %d, proof %d", len(pubkey), len(proof))
	}

	id := sdktypes.ConsAddress(common.HexToAddress(testValidator).Bytes())
	if err := types.VerifyRotationProof(testChainID, id, types.KeyTypeMlDsa65, pubkey, proof); err != nil {
		t.Fatalf("printed proof does not verify: %v", err)
	}

	// the file it wrote has to be one a node can actually sign with
	pv := privval.LoadFilePVEmptyState(out, "")
	if !bytes.Equal(pv.Key.PrivKey.PubKey().Bytes(), pubkey) {
		t.Fatal("the key written is not the key the proof names")
	}
}

// A proof is bound to the chain and to the validator, so neither can be
// swapped after the fact.
func TestRotateProofIsBoundToChainAndValidator(t *testing.T) {
	out := filepath.Join(t.TempDir(), "key.json")
	fields := runRotateProof(t,
		"--chain-id", testChainID, "--validator", testValidator, "--out", out)
	pubkey, proof := decode(t, fields["pubkey"]), decode(t, fields["proof"])

	id := sdktypes.ConsAddress(common.HexToAddress(testValidator).Bytes())
	other := sdktypes.ConsAddress(common.HexToAddress("0x000000000000000000000000000000000000dEaD").Bytes())

	if err := types.VerifyRotationProof("goat-testnet3", id, types.KeyTypeMlDsa65, pubkey, proof); err == nil {
		t.Fatal("proof replayed onto another chain")
	}
	if err := types.VerifyRotationProof(testChainID, other, types.KeyTypeMlDsa65, pubkey, proof); err == nil {
		t.Fatal("proof replayed onto another validator")
	}
}

// A consensus key is not something to clobber by accident.
func TestRotateProofRefusesToOverwriteAKey(t *testing.T) {
	out := filepath.Join(t.TempDir(), "key.json")
	runRotateProof(t, "--chain-id", testChainID, "--validator", testValidator, "--out", out)

	cmd := RotateProofCommand()
	cmd.SetOut(new(bytes.Buffer))
	cmd.SetErr(new(bytes.Buffer))
	cmd.SetArgs([]string{"--chain-id", testChainID, "--validator", testValidator, "--out", out})
	if err := cmd.Execute(); err == nil {
		t.Fatal("overwrote an existing key file")
	}
}

// --from-key signs with a key the operator already holds and writes nothing.
func TestRotateProofSignsWithAnExistingKey(t *testing.T) {
	out := filepath.Join(t.TempDir(), "key.json")
	first := runRotateProof(t, "--chain-id", testChainID, "--validator", testValidator, "--out", out)
	second := runRotateProof(t, "--chain-id", testChainID, "--validator", testValidator, "--from-key", out)

	if first["pubkey"] != second["pubkey"] || first["proof"] != second["proof"] {
		t.Fatal("--from-key did not reproduce the proof for the same key")
	}
	if _, ok := second["key"]; ok {
		t.Fatal("--from-key wrote a key file")
	}
}

// secp256k1 stays reachable: rotation is also the answer to a compromised key,
// not only to the post-quantum migration.
func TestRotateProofSupportsSecp256k1(t *testing.T) {
	out := filepath.Join(t.TempDir(), "key.json")
	fields := runRotateProof(t, "--chain-id", testChainID, "--validator", testValidator,
		"--key-type", types.KeyTypeSecp256k1, "--out", out)

	if fields["key_type"] != "0" {
		t.Fatalf("wire key type %q, want 0", fields["key_type"])
	}
	pubkey, proof := decode(t, fields["pubkey"]), decode(t, fields["proof"])
	id := sdktypes.ConsAddress(common.HexToAddress(testValidator).Bytes())
	if err := types.VerifyRotationProof(testChainID, id, types.KeyTypeSecp256k1, pubkey, proof); err != nil {
		t.Fatalf("secp256k1 proof does not verify: %v", err)
	}
}
