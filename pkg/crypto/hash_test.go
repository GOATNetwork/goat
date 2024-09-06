package crypto

import (
	"reflect"
	"testing"
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
