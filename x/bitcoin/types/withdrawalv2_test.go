package types

import (
	"encoding/hex"
	"testing"
)

func TestCalTxPrice(t *testing.T) {
	type fields struct {
		NewNoWitnessTx string
		NewTxFee       uint64
		WitnessSize    uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "43b516d1caeb13660cef2986a0b7edc238bd9ad8fcaef16baef03038ac4d5d21",
			fields: fields{
				NewNoWitnessTx: "020000000150b03d3cad9c662b29937c864db92d67c401a58b7bc973f00ffe4810c0cebf040000000000fdffffff027a72bf0100000000160014192e80ed2c7c412bdc2a6c8f371d15cb90f3c85b46bd0500000000001600147a970204fc2f2fbad75ebd83487c56a0af50a86600000000",
				NewTxFee:       5640,
				WitnessSize:    109,
			},
			want: 40.213903743315505,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := hex.DecodeString(tt.fields.NewNoWitnessTx)
			if err != nil {
				t.Errorf("hex.DecodeString() error = %v", err)
				return
			}

			withdrawal := &MsgProcessWithdrawalV2{
				NoWitnessTx: tx,
				TxFee:       tt.fields.NewTxFee,
				WitnessSize: tt.fields.WitnessSize,
			}
			if got := withdrawal.CalTxPrice(); got != tt.want {
				t.Errorf("MsgProcessWithdrawalV2.CalTxPrice() = %v, want %v", got, tt.want)
			}

			replace := &MsgReplaceWithdrawalV2{
				NewNoWitnessTx: tx,
				NewTxFee:       tt.fields.NewTxFee,
				WitnessSize:    tt.fields.WitnessSize,
			}
			if got := replace.CalTxPrice(); got != tt.want {
				t.Errorf("MsgReplaceWithdrawalV2.CalTxPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}
