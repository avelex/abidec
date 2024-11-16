package abidec

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	structNameRegex = regexp.MustCompile(`struct (\w+){?`)
	fieldRegex      = regexp.MustCompile(`(\w+(?:\[\d*\])?)\s+(\w+)`)

	getStructTemplate = `{"inputs": [], "name": "%s", "outputs": [{"components": [%s],"internalType": "struct %s", "name": "%s", "type": "tuple"}], "stateMutability": "pure", "type": "function"}`
)

type StructDef struct {
	Name          string
	Fields        []StructField
	InitialStrDef string
}

type StructField struct {
	Name string
	Type string
}

func (s StructDef) Getter() string {
	return "get" + s.Name
}

func StructDefFromString(structDef string) (StructDef, error) {
	structDef = strings.TrimSpace(structDef)

	nameMatch := structNameRegex.FindStringSubmatch(structDef)
	if len(nameMatch) != 2 {
		return StructDef{}, fmt.Errorf("invalid struct definition format")
	}
	structName := nameMatch[1]

	lines := strings.Split(structDef, "\n")
	fields := make([]StructField, 0)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "struct") || strings.HasPrefix(line, "{") || strings.HasPrefix(line, "}") {
			continue
		}

		match := fieldRegex.FindStringSubmatch(line)
		if len(match) == 3 {
			fields = append(fields, StructField{
				Type: match[1],
				Name: match[2],
			})
		}
	}

	if len(fields) == 0 {
		return StructDef{}, fmt.Errorf("no fields found in struct definition")
	}

	return StructDef{
		Name:          structName,
		Fields:        fields,
		InitialStrDef: structDef,
	}, nil
}

func CreateGetterMethodForStruct(s StructDef) (string, error) {
	if len(s.Fields) == 0 {
		return "", fmt.Errorf("struct has no fields")
	}

	components := make([]string, len(s.Fields))
	for i, field := range s.Fields {
		components[i] = fmt.Sprintf(`{"internalType": "%s", "name": "%s", "type": "%s"}`, field.Type, field.Name, field.Type)
	}

	return fmt.Sprintf(getStructTemplate, s.Getter(), strings.Join(components, ",\n"), s.Name, s.Name), nil
}
