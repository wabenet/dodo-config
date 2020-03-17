package types

import (
	"reflect"
	"strings"

	"github.com/dodo/dodo-config/pkg/decoder"
	"github.com/oclaussen/dodo/pkg/types"
)

func NewVolume() decoder.Producer {
	return func() (interface{}, decoder.Decoder) {
		target := &types.Volume{}
		return &target, Volume(&target)
	}
}

func Volume(target interface{}) decoder.Decoder {
	// TODO: wtf this cast
	vol := *(target.(**types.Volume))
	return decoder.Kinds(map[reflect.Kind]decoder.Decoder{
		reflect.Map: decoder.Keys(map[string]decoder.Decoder{
			"source":    decoder.String(&vol.Source),
			"target":    decoder.String(&vol.Target),
			"read_only": decoder.Bool(&vol.Readonly),
		}),
		reflect.String: func(s *decoder.Status, config interface{}) {
			var decoded string
			decoder.String(&decoded)(s, config)
			switch values := strings.SplitN(decoded, ":", 3); len(values) {
			case 1:
				vol.Source = values[0]
			case 2:
				vol.Source = values[0]
				vol.Target = values[1]
			case 3:
				vol.Source = values[0]
				vol.Target = values[1]
				vol.Readonly = values[2] == "ro"
			default:
				s.Error("invalid volume")
                                return
			}
		},
	})
}
