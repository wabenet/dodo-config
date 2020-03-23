package types

import (
	"reflect"

	"github.com/dodo/dodo-config/pkg/decoder"
	"github.com/oclaussen/dodo/pkg/types"
)

func NewBackdrop() decoder.Producer {
	return func() (interface{}, decoder.Decoding) {
		target := &types.Backdrop{Entrypoint: &types.Entrypoint{}}
		return &target, DecodeBackdrop(&target)
	}
}

func DecodeBackdrop(target interface{}) decoder.Decoding {
	// TODO: wtf this cast
	backdrop := *(target.(**types.Backdrop))
	return decoder.Kinds(map[reflect.Kind]decoder.Decoding{
		reflect.Map: decoder.Keys(map[string]decoder.Decoding{
			"name":           decoder.String(&backdrop.ContainerName),
			"alias":          decoder.Slice(decoder.NewString(), &backdrop.Aliases),
			"aliases":        decoder.Slice(decoder.NewString(), &backdrop.Aliases),
			"container_name": decoder.String(&backdrop.ContainerName),
			"image":          decoder.String(&backdrop.ImageId),
			"interactive":    decoder.Bool(&backdrop.Entrypoint.Interactive),
			"script":         decoder.String(&backdrop.Entrypoint.Script),
			"interpreter": decoder.Kinds(map[reflect.Kind]decoder.Decoding{
				reflect.String: decoder.Singleton(decoder.NewString(), &backdrop.Entrypoint.Interpreter),
				reflect.Slice:  decoder.Slice(decoder.NewString(), &backdrop.Entrypoint.Interpreter),
			}),
			"arguments": decoder.Kinds(map[reflect.Kind]decoder.Decoding{
				reflect.String: decoder.Singleton(decoder.NewString(), &backdrop.Entrypoint.Arguments),
				reflect.Slice:  decoder.Slice(decoder.NewString(), &backdrop.Entrypoint.Arguments),
			}),
			"command": decoder.Kinds(map[reflect.Kind]decoder.Decoding{
				reflect.String: decoder.Singleton(decoder.NewString(), &backdrop.Entrypoint.Arguments),
				reflect.Slice:  decoder.Slice(decoder.NewString(), &backdrop.Entrypoint.Arguments),
			}),
			"env": decoder.Kinds(map[reflect.Kind]decoder.Decoding{
				reflect.String: decoder.Singleton(NewEnvironment(), &backdrop.Environment),
				reflect.Slice:  decoder.Slice(NewEnvironment(), &backdrop.Environment),
			}),
			"environment": decoder.Kinds(map[reflect.Kind]decoder.Decoding{
				reflect.String: decoder.Singleton(NewEnvironment(), &backdrop.Environment),
				reflect.Slice:  decoder.Slice(NewEnvironment(), &backdrop.Environment),
			}),
			"volume": decoder.Kinds(map[reflect.Kind]decoder.Decoding{
				reflect.String: decoder.Singleton(NewVolume(), &backdrop.Volumes),
				reflect.Map:    decoder.Singleton(NewVolume(), &backdrop.Volumes),
				reflect.Slice:  decoder.Slice(NewVolume(), &backdrop.Volumes),
			}),
			"volumes": decoder.Kinds(map[reflect.Kind]decoder.Decoding{
				reflect.String: decoder.Singleton(NewVolume(), &backdrop.Volumes),
				reflect.Map:    decoder.Singleton(NewVolume(), &backdrop.Volumes),
				reflect.Slice:  decoder.Slice(NewVolume(), &backdrop.Volumes),
			}),
			"device": decoder.Kinds(map[reflect.Kind]decoder.Decoding{
				reflect.String: decoder.Singleton(NewDevice(), &backdrop.Devices),
				reflect.Map:    decoder.Singleton(NewDevice(), &backdrop.Devices),
				reflect.Slice:  decoder.Slice(NewDevice(), &backdrop.Devices),
			}),
			"devices": decoder.Kinds(map[reflect.Kind]decoder.Decoding{
				reflect.String: decoder.Singleton(NewDevice(), &backdrop.Devices),
				reflect.Map:    decoder.Singleton(NewDevice(), &backdrop.Devices),
				reflect.Slice:  decoder.Slice(NewDevice(), &backdrop.Devices),
			}),
			"ports": decoder.Kinds(map[reflect.Kind]decoder.Decoding{
				reflect.String: decoder.Singleton(NewPort(), &backdrop.Ports),
				reflect.Map:    decoder.Singleton(NewPort(), &backdrop.Ports),
				reflect.Slice:  decoder.Slice(NewPort(), &backdrop.Ports),
			}),
			"user":        decoder.String(&backdrop.User),
			"workdir":     decoder.String(&backdrop.WorkingDir),
			"working_dir": decoder.String(&backdrop.WorkingDir),
		}),
	})
}
