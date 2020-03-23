package types

import (
	"reflect"
	"strings"

	"github.com/dodo/dodo-config/pkg/decoder"
	"github.com/oclaussen/dodo/pkg/types"
)

func NewDevice() decoder.Producer {
	return func() (interface{}, decoder.Decoding) {
		target := &types.Device{}
		return &target, DecodeDevice(&target)
	}
}

func DecodeDevice(target interface{}) decoder.Decoding {
	// TODO: wtf this cast
	dev := *(target.(**types.Device))
	return decoder.Kinds(map[reflect.Kind]decoder.Decoding{
		reflect.Map: decoder.Keys(map[string]decoder.Decoding{
			"cgroup_rule": decoder.String(&dev.CgroupRule),
			"source":      decoder.String(&dev.Source),
			"target":      decoder.String(&dev.Target),
			"permissions": decoder.String(&dev.Permissions),
		}),
		reflect.String: func(d *decoder.Decoder, config interface{}) {
			var decoded string
			decoder.String(&decoded)(d, config)
			switch values := strings.SplitN(decoded, ":", 3); len(values) {
			case 1:
				dev.Source = values[0]
			case 2:
				dev.Source = values[0]
				dev.Target = values[1]
			case 3:
				dev.Source = values[0]
				dev.Target = values[1]
				dev.Permissions = values[2]
			default:
				d.Error("invalid device")
                                return
			}
		},
	})
}
