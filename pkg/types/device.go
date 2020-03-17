package types

import (
	"reflect"
	"strings"

	"github.com/dodo/dodo-config/pkg/decoder"
	"github.com/oclaussen/dodo/pkg/types"
)

func NewDevice() decoder.Producer {
	return func() (interface{}, decoder.Decoder) {
		target := &types.Device{}
		return &target, Device(&target)
	}
}

func Device(target interface{}) decoder.Decoder {
	// TODO: wtf this cast
	dev := *(target.(**types.Device))
	return decoder.Kinds(map[reflect.Kind]decoder.Decoder{
		reflect.Map: decoder.Keys(map[string]decoder.Decoder{
			"cgroup_rule": decoder.String(&dev.CgroupRule),
			"source":      decoder.String(&dev.Source),
			"target":      decoder.String(&dev.Target),
			"permissions": decoder.String(&dev.Permissions),
		}),
		reflect.String: func(s *decoder.Status, config interface{}) {
			var decoded string
			decoder.String(&decoded)(s, config)
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
				s.Error("invalid device")
                                return
			}
		},
	})
}
