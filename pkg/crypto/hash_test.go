package crypto

import (
	"crypto/rand"
	"fmt"
	"reflect"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
)

func TestUint64LE(t *testing.T) {
	type args struct {
		n []uint64
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "1",
			args: args{[]uint64{100, 1e4}},
			want: []byte{100, 0, 0, 0, 0, 0, 0, 0, 16, 39, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "2",
			args: args{[]uint64{4294967297}},
			want: []byte{1, 0, 0, 0, 1, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Uint64LE(tt.args.n...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uint64LE() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHash160Sum(t *testing.T) {
	t.Parallel()
	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("Hash160-%d", i), func(t *testing.T) {
			data := make([]byte, 32)
			_, _ = rand.Read(data)
			if got, want := Hash160Sum(data), btcutil.Hash160(data); !reflect.DeepEqual(got, want) {
				t.Errorf("Hash160Sum() = %x, want %x", got, data)
			}
		})
	}
}
