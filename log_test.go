package abidec_test

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/avelex/abidec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

func Test_Event(t *testing.T) {
	testCases := []struct {
		desc string
		sig  string
		log  *types.Log
		want map[string]any
	}{
		{
			desc: "Valid event",
			sig:  "event Transfer(address indexed from, address indexed to, uint256 value)",
			log: &types.Log{
				Topics: []common.Hash{
					common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
					common.HexToHash("0x000000000000000000000000f5213a6a2f0890321712520b8048d9886c1a9900"),
					common.HexToHash("0x0000000000000000000000006ac6b053a2858bea8ad758db680198c16e523184"),
				},
				Data: hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000001ec0ef243"),
			},
			want: map[string]any{
				"from":  common.HexToAddress("0xf5213a6a2f0890321712520b8048D9886c1A9900"),
				"to":    common.HexToAddress("0x6ac6b053a2858bea8ad758db680198c16e523184"),
				"value": big.NewInt(8255369795),
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ev, err := abidec.ParseEventSignature(tC.sig)
			if err != nil {
				t.Fatalf("failed to parse event signature: %v", err)
			}

			out := make(map[string]any)
			if err := abidec.ParseLogIntoMap(ev, out, tC.log); err != nil {
				t.Fatalf("failed to parse log: %v", err)
			}

			if !reflect.DeepEqual(out, tC.want) {
				t.Fatalf("expected %v, got %v", tC.want, out)
			}
		})
	}
}
