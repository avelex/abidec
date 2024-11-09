package abidec_test

import (
	"testing"

	"github.com/avelex/abidec"
)

func Test_StructDefFromString(t *testing.T) {
	testCases := []struct {
		desc      string
		structDef string
		want      abidec.StructDef
		wantErr   bool
	}{
		{
			desc: "Valid struct",
			structDef: `
			struct Task {
				string title;
				string description;
				address reporter;
				address assignee;
				uint256[2] deadline;
			}
			`,
			want: abidec.StructDef{
				Name: "Task",
				Fields: []abidec.StructField{
					{Name: "title", Type: "string"},
					{Name: "description", Type: "string"},
					{Name: "reporter", Type: "address"},
					{Name: "assignee", Type: "address"},
					{Name: "deadline", Type: "uint256[2]"},
				},
			},
			wantErr: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := abidec.StructDefFromString(tC.structDef)
			if tC.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !tC.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Name != tC.want.Name {
				t.Fatalf("expected name %s, got %s", tC.want.Name, got.Name)
			}

			if len(got.Fields) != len(tC.want.Fields) {
				t.Fatalf("expected %d fields, got %d", len(tC.want.Fields), len(got.Fields))
			}

			for i, field := range got.Fields {
				if field.Name != tC.want.Fields[i].Name {
					t.Fatalf("expected name %s, got %s", tC.want.Fields[i].Name, field.Name)
				}
			}
		})
	}
}

func Test_CreateGetterMethodForStruct(t *testing.T) {
	testCases := []struct {
		desc      string
		structDef abidec.StructDef
		want      string
		wantErr   bool
	}{
		{
			desc: "Valid struct",
			structDef: abidec.StructDef{
				Name: "Task",
				Fields: []abidec.StructField{
					{Name: "title", Type: "string"},
					{Name: "description", Type: "string"},
				},
			},
			want: `{"inputs": [], "name": "getTask", "outputs": [{"components": [{"internalType": "string", "name": "title", "type": "string"},
{"internalType": "string", "name": "description", "type": "string"}],"internalType": "struct Task", "name": "Task", "type": "tuple"}], "stateMutability": "pure", "type": "function"}`,
			wantErr: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := abidec.CreateGetterMethodForStruct(tC.structDef)
			if tC.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !tC.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tC.want {
				t.Fatalf("expected %s, got %s", tC.want, got)
			}
		})
	}
}
