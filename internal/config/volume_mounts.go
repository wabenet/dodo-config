package config

import (
	"cuelang.org/go/cue"
	"github.com/dodo-cli/dodo-config/pkg/cuetils"
	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/config"
	"github.com/hashicorp/go-multierror"
)

func VolumeMountsFromValue(v cue.Value) ([]*api.VolumeMount, error) {
	var errs error

	if out, err := VolumeMountsFromMap(v); err == nil {
		return out, err
	} else {
		errs = multierror.Append(errs, err)
	}

	if out, err := VolumeMountsFromList(v); err == nil {
		return out, err
	} else {
		errs = multierror.Append(errs, err)
	}

	return nil, errs
}

func VolumeMountsFromMap(v cue.Value) ([]*api.VolumeMount, error) {
	out := []*api.VolumeMount{}

	err := cuetils.IterMap(v, func(name string, v cue.Value) error {
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

	err := cuetils.IterList(v, func(v cue.Value) error {
		r, err := VolumeMountFromValue("", v)
		if err == nil {
			out = append(out, r)
		}

		return err
	})

	return out, err
}

func VolumeMountFromValue(name string, v cue.Value) (*api.VolumeMount, error) {
	var errs error

	if out, err := VolumeMountFromString(name, v); err == nil {
		return out, err
	} else {
		errs = multierror.Append(errs, err)
	}

	if out, err := VolumeMountFromStruct(name, v); err == nil {
		return out, err
	} else {
		errs = multierror.Append(errs, err)
	}

	return nil, errs
}

func VolumeMountFromString(_ string, v cue.Value) (*api.VolumeMount, error) {
	s, err := StringFromValue(v)
	if err != nil {
		return nil, err
	}

	return config.ParseVolumeMount(s)
}

func VolumeMountFromStruct(name string, v cue.Value) (*api.VolumeMount, error) {
	out := &api.VolumeMount{Source: name}

	if p, ok := cuetils.Get(v, "source"); ok {
		if v, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.Source = v
		}
	}

	if p, ok := cuetils.Get(v, "target"); ok {
		if v, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.Target = v
		}
	}

	if p, ok := cuetils.Get(v, "readonly"); ok {
		if v, err := p.Bool(); err != nil {
			return nil, err
		} else {
			out.Readonly = v
		}
	}

	return out, nil
}
