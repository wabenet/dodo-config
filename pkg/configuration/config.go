package configuration

import (
	"reflect"

	"github.com/dodo/dodo-config/pkg/decoder"
	cfgtypes "github.com/dodo/dodo-config/pkg/types"
	"github.com/oclaussen/dodo/pkg/types"
)

type Config struct {
	Backdrops map[string]*types.Backdrop
	Groups    map[string]*Config
}

func NewConfig() (interface{}, decoder.Decoder) {
	target := &Config{Groups: map[string]*Config{}, Backdrops: map[string]*types.Backdrop{}}
	return &target, decodeConfig(&target)
}

func decodeConfig(target interface{}) decoder.Decoder {
	// TODO: wtf this cast
	group := *(target.(**Config))
	dummy := []interface{}{}
	return decoder.Keys(map[string]decoder.Decoder{
		"groups":    decoder.Map(NewConfig, &group.Groups),
		"backdrops": decoder.Map(cfgtypes.NewBackdrop(), &group.Backdrops),
		"include": decoder.Kinds(map[reflect.Kind]decoder.Decoder{
			reflect.Map:   decoder.Singleton(newInclude(target), &dummy),
			reflect.Slice: decoder.Slice(newInclude(target), &dummy),
		}),
	})
}

func newInclude(target interface{}) decoder.Producer {
	return func() (interface{}, decoder.Decoder) {
		var dummy struct{}
		return &dummy, decoder.Keys(map[string]decoder.Decoder{
			"text": decoder.IncludeText(decodeConfig(target)),
			"file": decoder.IncludeFile(decodeConfig(target)),
		})
	}
}
