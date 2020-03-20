package types

import (
	"reflect"
	"strings"

	"github.com/dodo/dodo-config/pkg/decoder"
	"github.com/oclaussen/dodo/pkg/types"
)

func NewPort() decoder.Producer {
	return func() (interface{}, decoder.Decoder) {
		target := &types.Port{}
		return &target, Port(&target)
	}
}

func Port(target interface{}) decoder.Decoder {
	// TODO: wtf this cast
	port := *(target.(**types.Port))
	return decoder.Kinds(map[reflect.Kind]decoder.Decoder{
		reflect.Map: decoder.Keys(map[string]decoder.Decoder{
			"target":    decoder.String(&port.Target),
			"published": decoder.String(&port.Published),
			"protocol":  decoder.String(&port.Protocol),
			"host_ip":   decoder.String(&port.HostIp),
		}),
		reflect.String: func(s *decoder.Status, config interface{}) {
			var decoded string
			decoder.String(&decoded)(s, config)
			switch values := strings.SplitN(decoded, ":", 3); len(values) {
			case 1:
				port.Target = values[0]
			case 2:
				port.Published = values[0]
				port.Target = values[1]
			case 3:
				port.HostIp = values[0]
				port.Published = values[1]
				port.Target = values[2]
			default:
				s.Error("invalid port definition")
                                return
			}
			switch values := strings.SplitN(port.Target, "/", 2); len(values) {
			case 1:
				port.Target = values[0]
			case 2:
				port.Target = values[0]
				port.Protocol = values[1]
			default:
				s.Error("invalid port definition")
                                return
			}
		},
	})
}