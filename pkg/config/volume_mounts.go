package config

import (
	"cuelang.org/go/cue"
	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/config"
)

func VolumeMountsFromValue(v cue.Value) ([]*api.VolumeMount, error) {
	if out, err := VolumeMountsFromMap(v); err == nil {
		return out, err
	}

	if out, err := VolumeMountsFromList(v); err == nil {
		return out, err
	}

	return nil, ErrUnexpectedSpec
}

func VolumeMountsFromMap(v cue.Value) ([]*api.VolumeMount, error) {
	out := []*api.VolumeMount{}

	err := eachInMap(v, func(name string, v cue.Value) error {
		r, err := VolumeMountFromValue(name, v)
		if err == nil {
			out = append(out, r)
		}

		return err

	})

	return out, err
}

func VolumeMountsFromList(v cue.Value) ([]*api.VolumeMount, error) {
	out := []*api.VolumeMount{}

	err := eachInList(v, func(v cue.Value) error {
		r, err := VolumeMountFromValue("", v)
		if err == nil {
			out = append(out, r)
		}

		return err
	})

	return out, err
}

func VolumeMountFromValue(name string, v cue.Value) (*api.VolumeMount, error) {
	if out, err := VolumeMountFromString(name, v); err == nil {
		return out, err
	}

	if out, err := VolumeMountFromStruct(name, v); err == nil {
		return out, err
	}

	return nil, ErrUnexpectedSpec
}

func VolumeMountFromString(_ string, v cue.Value) (*api.VolumeMount, error) {
	s, err := v.String()
	if err != nil {
		return nil, err
	}

	return config.ParseVolumeMount(s)
}

func VolumeMountFromStruct(name string, v cue.Value) (*api.VolumeMount, error) {
	out := &api.VolumeMount{Source: name}

	if p, ok := property(v, "source"); ok {
		if v, err := p.String(); err != nil {
			return nil, err
		} else {
			out.Source = v
		}
	}

	if p, ok := property(v, "target"); ok {
		if v, err := p.String(); err != nil {
			return nil, err
		} else {
			out.Target = v
		}
	}

	if p, ok := property(v, "readonly"); ok {
		if v, err := p.Bool(); err != nil {
			return nil, err
		} else {
			out.Readonly = v
		}
	}

	return out, nil
}
