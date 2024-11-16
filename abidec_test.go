package abidec_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"strings"
	"testing"

	"github.com/avelex/abidec"
)

func Test_AbiDecoder_Struct_Decode(t *testing.T) {
	testCases := []struct {
		desc       string
		structDef  string
		params     string
		wantStruct string
		wantValues map[string]any
		wantError  bool
	}{
		{
			desc: "Valid struct",
			structDef: `
			struct Task {
				string title;
				string description;
				address reporter;
				address assignee
				uint256[2] deadline;
			}
			`,
			params:     "000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000100000000000000000000000000acdad15d8f07d8df258fe11332b752785d6b1d220000000000000000000000008b1383709d1e80a291de5d67993252dfc52c370000000000000000000000000000000000000000000000000000000000672f269b00000000000000000000000000000000000000000000000000000000672f34ab0000000000000000000000000000000000000000000000000000000000000006526f636b65740000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001943726561746520526f636b657420546f20546865204d6f6f6e00000000000000",
			wantStruct: "Task",
			wantValues: map[string]any{
				"Task": map[string]any{
					"title":       "Rocket",
					"description": "Create Rocket To The Moon",
					"reporter":    "0xaCDaD15d8F07D8Df258fe11332b752785d6b1d22",
					"assignee":    "0x8B1383709D1e80A291DE5d67993252dFC52C3700",
					"deadline":    [2]uint64{1731143323, 1731146923},
				},
			},
			wantError: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			decoder, err := abidec.NewAbiDecoder(
				abidec.WithStruct(tC.structDef),
			)
			if err != nil {
				t.Fatalf("failed to create decoder: %v", err)
			}

			paramsBytes, err := hex.DecodeString(tC.params)
			if err != nil {
				t.Fatalf("failed to decode params: %v", err)
			}

			gotValues, err := decoder.Decode(paramsBytes)
			if err != nil {
				t.Fatalf("failed to decode params: %v", err)
			}

			gotValuesBytes, err := json.Marshal(&gotValues)
			if err != nil {
				t.Fatalf("failed to marshal gotValues: %v", err)
			}

			wantValuesBytes, err := json.Marshal(&tC.wantValues)
			if err != nil {
				t.Fatalf("failed to marshal wantValues: %v", err)
			}

			if !bytes.Equal(gotValuesBytes, wantValuesBytes) {
				t.Errorf("got: %s, want: %s", string(gotValuesBytes), string(wantValuesBytes))
			}
		})
	}
}

func Test_AbiDecoder_Struct_DecodeStruct(t *testing.T) {
	type TestStruct struct {
		Title       string
		Description string
		Reporter    [20]byte
		Assignee    [20]byte
		Deadline    [2]*big.Int
	}

	testCases := []struct {
		desc       string
		structDef  string
		params     string
		wantStruct TestStruct
		wantError  bool
	}{
		{
			desc: "Valid struct",
			structDef: `
			struct Task {
				string title;
				string description;
				address reporter;
				address assignee
				uint256[2] deadline;
			}
			`,
			params: "000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000100000000000000000000000000acdad15d8f07d8df258fe11332b752785d6b1d220000000000000000000000008b1383709d1e80a291de5d67993252dfc52c370000000000000000000000000000000000000000000000000000000000672f269b00000000000000000000000000000000000000000000000000000000672f34ab0000000000000000000000000000000000000000000000000000000000000006526f636b65740000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001943726561746520526f636b657420546f20546865204d6f6f6e00000000000000",
			wantStruct: TestStruct{
				Title:       "Rocket",
				Description: "Create Rocket To The Moon",
				Reporter:    [20]byte(mustHexDecode("0xaCDaD15d8F07D8Df258fe11332b752785d6b1d22")),
				Assignee:    [20]byte(mustHexDecode("0x8B1383709D1e80A291DE5d67993252dFC52C3700")),
				Deadline:    [2]*big.Int{big.NewInt(1731143323), big.NewInt(1731146923)},
			},
			wantError: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			decoder, err := abidec.NewAbiDecoder(
				abidec.WithStruct(tC.structDef),
			)
			if err != nil {
				t.Fatalf("failed to create decoder: %v", err)
			}

			paramsBytes, err := hex.DecodeString(tC.params)
			if err != nil {
				t.Fatalf("failed to decode params: %v", err)
			}

			var gotStruct TestStruct
			if err := decoder.DecodeStruct(paramsBytes, &gotStruct); err != nil {
				t.Fatalf("failed to decode params to struct: %v", err)
			}

			if gotStruct.Title != tC.wantStruct.Title {
				t.Fatalf("got title: %s, want: %s", gotStruct.Title, tC.wantStruct.Title)
			}

			if gotStruct.Description != tC.wantStruct.Description {
				t.Fatalf("got description: %s, want: %s", gotStruct.Description, tC.wantStruct.Description)
			}

			if !bytes.Equal(gotStruct.Reporter[:], tC.wantStruct.Reporter[:]) {
				t.Fatalf("got reporter: %x, want: %x", gotStruct.Reporter[:], tC.wantStruct.Reporter[:])
			}

			if !bytes.Equal(gotStruct.Assignee[:], tC.wantStruct.Assignee[:]) {
				t.Fatalf("got assignee: %x, want: %x", gotStruct.Assignee[:], tC.wantStruct.Assignee[:])
			}

			if !bytes.Equal(gotStruct.Deadline[0].Bytes(), tC.wantStruct.Deadline[0].Bytes()) {
				t.Fatalf("got deadline[0]: %x, want: %x", gotStruct.Deadline[0].Bytes(), tC.wantStruct.Deadline[0].Bytes())
			}

			if !bytes.Equal(gotStruct.Deadline[1].Bytes(), tC.wantStruct.Deadline[1].Bytes()) {
				t.Fatalf("got deadline[1]: %x, want: %x", gotStruct.Deadline[1].Bytes(), tC.wantStruct.Deadline[1].Bytes())
			}
		})
	}
}

func mustHexDecode(s string) []byte {
	b, err := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	if err != nil {
		panic(err)
	}
	return b
}
