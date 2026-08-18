package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	cmttypes "github.com/cometbft/cometbft/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
)

// The shipped genesis files are consensus critical: a node joining the
// network parses them with whatever cometbft and sdk versions the binary was
// built against. A dependency bump that tightens genesis or consensus param
// validation would only surface when someone tries to sync from scratch, so
// parse and validate them here instead.
func TestShippedGenesisParses(t *testing.T) {
	entries, err := genesisFiles.ReadDir("genesis")
	if err != nil {
		t.Fatalf("read embedded genesis dir: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("no embedded genesis files")
	}

	for _, entry := range entries {
		t.Run(entry.Name(), func(t *testing.T) {
			raw, err := genesisFiles.ReadFile(filepath.Join("genesis", entry.Name()))
			if err != nil {
				t.Fatalf("read: %v", err)
			}

			path := filepath.Join(t.TempDir(), entry.Name())
			if err := os.WriteFile(path, raw, 0o600); err != nil {
				t.Fatalf("write: %v", err)
			}

			appGenesis, err := genutiltypes.AppGenesisFromFile(path)
			if err != nil {
				t.Fatalf("AppGenesisFromFile: %v", err)
			}
			if err := appGenesis.ValidateAndComplete(); err != nil {
				t.Fatalf("ValidateAndComplete: %v", err)
			}

			// the consensus params are what a dependency bump is most likely
			// to start rejecting, pub_key_types above all
			doc, err := appGenesis.ToGenesisDoc()
			if err != nil {
				t.Fatalf("ToGenesisDoc: %v", err)
			}
			if err := doc.ValidateAndComplete(); err != nil {
				t.Fatalf("GenesisDoc.ValidateAndComplete: %v", err)
			}
			if err := doc.ConsensusParams.ValidateBasic(); err != nil {
				t.Fatalf("ConsensusParams.ValidateBasic: %v", err)
			}
			for _, kt := range doc.ConsensusParams.Validator.PubKeyTypes {
				if !cmttypes.IsValidPubkeyType(doc.ConsensusParams.Validator, kt) {
					t.Fatalf("pub_key_type %q rejected by this cometbft version", kt)
				}
			}

			var appState map[string]json.RawMessage
			if err := json.Unmarshal(appGenesis.AppState, &appState); err != nil {
				t.Fatalf("app_state: %v", err)
			}
			for _, module := range []string{"auth", "locking", "goat", "relayer", "bitcoin"} {
				if _, ok := appState[module]; !ok {
					t.Errorf("app_state is missing the %s module", module)
				}
			}
		})
	}
}
