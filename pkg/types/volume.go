package types

import (
	"reflect"
	"strings"

	"github.com/dodo/dodo-config/pkg/decoder"
	"github.com/oclaussen/dodo/pkg/types"
)

func NewVolume() decoder.Producer {
	return func() (interface{}, decoder.Decoding) {
		target := &types.Volume{}
		return &target, DecodeVolume(&target)
	}
}

func DecodeVolume(target interface{}) decoder.Decoding {
	// TODO: wtf this cast
	vol := *(target.(**types.Volume))
	return decoder.Kinds(map[reflect.Kind]decoder.Decoding{
		reflect.Map: decoder.Keys(map[string]decoder.Decoding{
			"source":    decoder.String(&vol.Source),
			"target":    decoder.String(&vol.Target),
			"read_only": decoder.Bool(&vol.Readonly),
		}),
		reflect.String: func(d *decoder.Decoder, config interface{}) {
			var decoded string
			decoder.String(&decoded)(d, config)
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
				d.Error("invalid volume")
                                return
			}
		},
	})
}
