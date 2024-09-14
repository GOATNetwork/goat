package types

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/ethereum/go-ethereum/common"
)

func TestVerifyMerkelProof(t *testing.T) {
	root1, err := chainhash.NewHashFromStr(
		"6e026bd12f312dce9ad309cad9998a18d560a54467ce74491d1997c316637763")
	if err != nil {
		panic(err)
	}

	root2, err := chainhash.NewHashFromStr(
		"59443731d81f2171c6bcbd7019e6762d87adf06ce93cb3a9ad6c699350040402")
	if err != nil {
		panic(err)
	}

	type args struct {
		txid  string
		proof []byte
		root  []byte
		index uint32
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "height-1-0",
			args: args{
				txid:  "59443731d81f2171c6bcbd7019e6762d87adf06ce93cb3a9ad6c699350040402",
				root:  root2.CloneBytes(),
				proof: nil,
				index: 0,
			},
			want: true,
		},
		{
			name: "false",
			args: args{
				txid:  "6593c3e2836908ffbe9fa27238629dcf609baeef1c9a3521c1522aa56c163b37",
				root:  common.Hex2Bytes("00"),
				proof: nil,
				index: 0,
			},
			want: false,
		},
		{
			name: "height-101-0",
			args: args{
				txid:  "fed5dbea421b4c341a2abb94d90ad02a429fd012268c3b8fdd5d03433a8a189d",
				root:  root1.CloneBytes(),
				proof: common.Hex2Bytes("07ca9a98aacf23988c8f2ab214d0cdd55ca7c4e43b8adb9f6710ee941ec6233efc2f3ea9ed46f53d8a21952801138c4555c8e240e02540341fb496ce2d4db15515b9492e91f2e8bdd92981b0be3e76bcd08f05436811bf0203a49c232e3d3985"),
				index: 0,
			},
			want: true,
		},
		{
			name: "height-101-1",
			args: args{
				txid:  "3e23c61e94ee10679fdb8a3be4c4a75cd5cdd014b22a8f8c9823cfaa989aca07",
				root:  root1.CloneBytes(),
				proof: common.Hex2Bytes("9d188a3a43035ddd8f3b8c2612d09f422ad00ad994bb2a1a344c1b42eadbd5fefc2f3ea9ed46f53d8a21952801138c4555c8e240e02540341fb496ce2d4db15515b9492e91f2e8bdd92981b0be3e76bcd08f05436811bf0203a49c232e3d3985"),
				index: 1,
			},
			want: true,
		},
		{
			name: "height-101-2",
			args: args{
				txid:  "4cdbb27b088d397e1dbede283b6aaaa6cb3c723c779b5e331f8e6a8c41470793",
				root:  root1.CloneBytes(),
				proof: common.Hex2Bytes("776769f1d287cd2d8fbe0252935e185cc3f7d6eee2d555c5f6bd2acb30e6a4f3c00ae98a09a9f63c32c7a26f1aede1dc37c0800fff76f9c74fbd19b66f40164215b9492e91f2e8bdd92981b0be3e76bcd08f05436811bf0203a49c232e3d3985"),
				index: 2,
			},
			want: true,
		},
		{
			name: "height-101-3",
			args: args{
				txid:  "f3a4e630cb2abdf6c555d5e2eed6f7c35c185e935202be8f2dcd87d2f1696777",
				root:  root1.CloneBytes(),
				proof: common.Hex2Bytes("930747418c6a8e1f335e9b773c723ccba6aa6a3b28debe1d7e398d087bb2db4cc00ae98a09a9f63c32c7a26f1aede1dc37c0800fff76f9c74fbd19b66f40164215b9492e91f2e8bdd92981b0be3e76bcd08f05436811bf0203a49c232e3d3985"),
				index: 3,
			},
			want: true,
		},
		{
			name: "height-101-4",
			args: args{
				txid:  "3865df31da58bd53107f4aa90f131cfa60cf564739e3e6b6d8edda4a3047dcb0",
				root:  root1.CloneBytes(),
				proof: common.Hex2Bytes("f9478334ab22a8d89bcfd45a48552385a1022ae05111b514a574ec1ea62ab07cba810d97eccaa2674cd00fe26e47b04a11a97dcc2d5bb29c6d0bf68643c62fe24660b25ea76d43b2465c7ed53e49b637124863f1104d6f137ef1345491bba22e"),
				index: 4,
			},
			want: true,
		},
		{
			name: "height-101-5",
			args: args{
				txid:  "7cb02aa61eec74a514b51151e02a02a1852355485ad4cf9bd8a822ab348347f9",
				root:  root1.CloneBytes(),
				proof: common.Hex2Bytes("b0dc47304adaedd8b6e6e3394756cf60fa1c130fa94a7f1053bd58da31df6538ba810d97eccaa2674cd00fe26e47b04a11a97dcc2d5bb29c6d0bf68643c62fe24660b25ea76d43b2465c7ed53e49b637124863f1104d6f137ef1345491bba22e"),
				index: 5,
			},
			want: true,
		},
		{
			name: "height-101-6",
			args: args{
				txid:  "da4ba88617d61e5cfae541424ee5aa9f33e8677fe0b678311a0e7b82fc445ee1",
				root:  root1.CloneBytes(),
				proof: common.Hex2Bytes("e15e44fc827b0e1a3178b6e07f67e8339faae54e4241e5fa5c1ed61786a84bdae65ad2384e0316a78d5abc340b02a56c0e08d4af8e73fbb2de3526f829bc53944660b25ea76d43b2465c7ed53e49b637124863f1104d6f137ef1345491bba22e"),
				index: 6,
			},
			want: true,
		},
	}

	t.Parallel()
	for idx, tt := range tests {
		idx, tt := idx, tt
		t.Run(tt.name, func(t *testing.T) {
			txid, err := chainhash.NewHashFromStr(tt.args.txid)
			if err != nil {
				t.Errorf("invalid txid: %d", idx)
				return
			}

			if got := VerifyMerkelProof(txid[:], tt.args.root, tt.args.proof, tt.args.index); got != tt.want {
				t.Errorf("VerifyMerkelProof() = %v, want %v", got, tt.want)
			}
		})
	}
}
