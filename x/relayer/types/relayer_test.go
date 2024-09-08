package types

import (
	"crypto/sha256"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
)

func TestRelayer_Threshold(t *testing.T) {
	type fields struct {
		Voters []string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name:   "length 1",
			fields: fields{},
			want:   1,
		},
		{
			name:   "length 2",
			fields: fields{Voters: make([]string, 1)},
			want:   2,
		},
		{
			name:   "length 3",
			fields: fields{Voters: make([]string, 2)},
			want:   2,
		},
		{
			name:   "length 7",
			fields: fields{Voters: make([]string, 6)},
			want:   5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			relayer := &Relayer{Voters: tt.fields.Voters}
			if got := relayer.Threshold(); got != tt.want {
				t.Errorf("Relayer.Threshold() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVoteSignDoc(t *testing.T) {
	type args struct {
		method   string
		chainId  string
		proposer string
		sequence uint64
		epoch    uint64
		data     []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "1",
			args: args{
				method:   "1",
				chainId:  "goat-test-1",
				proposer: "goat15lftp68cgrca3t6j8w6s6j6v5qedfz2wkt7y9a",
				sequence: 100,
				epoch:    1,
				data:     hexutil.MustDecode("0x8f2dc0366a151063704578eca4ab971f18e801a7e70d40dbf9da80258cb88cbb"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasher := sha256.New()
			hasher.Write([]byte(tt.args.chainId))
			hasher.Write(goatcrypto.Uint64LE(tt.args.sequence, tt.args.epoch))
			hasher.Write([]byte(tt.args.method))
			hasher.Write([]byte(tt.args.proposer))
			hasher.Write(tt.args.data)
			want := hasher.Sum(nil)

			if got := VoteSignDoc(tt.args.method, tt.args.chainId, tt.args.proposer, tt.args.sequence, tt.args.epoch, tt.args.data); !reflect.DeepEqual(got, want) {
				t.Errorf("VoteSignDoc() = %v, want %v", got, tt.want)
			}
		})
	}
}
