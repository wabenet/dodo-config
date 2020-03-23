package types

import (
	"os"
	"strings"

	"github.com/dodo/dodo-config/pkg/decoder"
	"github.com/oclaussen/dodo/pkg/types"
)

func NewEnvironment() decoder.Producer {
	return func() (interface{}, decoder.Decoding) {
		target := &types.Environment{}
		return &target, DecodeEnvironment(&target)
	}
}

func DecodeEnvironment(target interface{}) decoder.Decoding {
	// TODO: wtf this cast
	env := *(target.(**types.Environment))
	return func(d *decoder.Decoder, config interface{}) {
		var decoded string
		decoder.String(&decoded)(d, config)
		switch values := strings.SplitN(decoded, "=", 2); len(values) {
		case 1:
			env.Key = values[0]
			env.Value = os.Getenv(values[0])
		case 2:
			env.Key = values[0]
			env.Value = values[1]
		default:
			d.Error("invalid environment")
		}
	}
}
