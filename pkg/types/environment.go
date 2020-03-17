package types

import (
	"os"
	"strings"

	"github.com/dodo/dodo-config/pkg/decoder"
	"github.com/oclaussen/dodo/pkg/types"
)

func NewEnvironment() decoder.Producer {
	return func() (interface{}, decoder.Decoder) {
		target := &types.Environment{}
		return &target, Environment(&target)
	}
}

func Environment(target interface{}) decoder.Decoder {
	// TODO: wtf this cast
	env := *(target.(**types.Environment))
	return func(s *decoder.Status, config interface{}) {
		var decoded string
		decoder.String(&decoded)(s, config)
		switch values := strings.SplitN(decoded, "=", 2); len(values) {
		case 1:
			env.Key = values[0]
			env.Value = os.Getenv(values[0])
		case 2:
			env.Key = values[0]
			env.Value = values[1]
		default:
			s.Error("invalid environment")
		}
	}
}
