package abidec

import (
	"fmt"

	"github.com/umbracle/ethgo/abi"
)

type DecoderOpts func(*AbiDecoder) error

func WithStruct(structDef string) DecoderOpts {
	return func(d *AbiDecoder) error {
		def, err := StructDefFromString(structDef)
		if err != nil {
			return fmt.Errorf("invalid struct definition: %w", err)
		}

		if err := d.buildAbiFromStructDef(def); err != nil {
			return fmt.Errorf("failed to build abi from struct definition: %w", err)
		}

		d.structDef = structDef

		return nil
	}
}

type AbiDecoder struct {
	abi       *abi.ABI
	method    string
	structDef string
}

func NewAbiDecoder(opts ...DecoderOpts) (*AbiDecoder, error) {
	decoder := &AbiDecoder{}
	for _, opt := range opts {
		if err := opt(decoder); err != nil {
			return nil, err
		}
	}
	return decoder, nil
}

func (d *AbiDecoder) Decode(data []byte) (map[string]any, error) {
	if d.abi == nil {
		return nil, fmt.Errorf("abi not initialized")
	}

	method := d.abi.GetMethod(d.method)
	if method == nil {
		return nil, fmt.Errorf("method not found")
	}
	return method.Decode(data)
}

func (d *AbiDecoder) buildAbiFromStructDef(s StructDef) error {
	method, err := CreateGetterMethodForStruct(s)
	if err != nil {
		return fmt.Errorf("failed to create getter method for struct: %w", err)
	}

	abiStr := "[" + method + "]"
	d.abi, err = abi.NewABI(abiStr)
	if err != nil {
		return fmt.Errorf("failed to create abi: %w", err)
	}

	d.method = s.Getter()

	return nil
}
