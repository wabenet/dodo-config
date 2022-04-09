package config

import (
	"cuelang.org/go/cue"
	"github.com/dodo-cli/dodo-config/pkg/cuetils"
	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/config"
	"github.com/hashicorp/go-multierror"
)

func DeviceMappingsFromValue(v cue.Value) ([]*api.DeviceMapping, error) {
	var errs error

	if out, err := DeviceMappingsFromMap(v); err == nil {
		return out, nil
	} else {
		errs = multierror.Append(errs, err)
	}

	if out, err := DeviceMappingsFromList(v); err == nil {
		return out, nil
	} else {
		errs = multierror.Append(errs, err)
	}

	return nil, errs
}

func DeviceMappingsFromMap(v cue.Value) ([]*api.DeviceMapping, error) {
	out := []*api.DeviceMapping{}

	err := cuetils.IterMap(v, func(name string, v cue.Value) error {
		r, err := DeviceMappingFromValue(name, v)
		if err == nil {
			out = append(out, r)
		}

		return err

	})

	return out, err
}

func DeviceMappingsFromList(v cue.Value) ([]*api.DeviceMapping, error) {
	out := []*api.DeviceMapping{}

	err := cuetils.IterList(v, func(v cue.Value) error {
		r, err := DeviceMappingFromValue("", v)
		if err == nil {
			out = append(out, r)
		}

		return err
	})

	return out, err
}

func DeviceMappingFromValue(name string, v cue.Value) (*api.DeviceMapping, error) {
	var errs error

	if out, err := DeviceMappingFromString(name, v); err == nil {
		return out, nil
	} else {
		errs = multierror.Append(errs, err)
	}

	if out, err := DeviceMappingFromStruct(name, v); err == nil {
		return out, nil
	} else {
		errs = multierror.Append(errs, err)
	}

	return nil, errs
}

func DeviceMappingFromString(_ string, v cue.Value) (*api.DeviceMapping, error) {
	s, err := StringFromValue(v)
	if err != nil {
		return nil, err
	}

	return config.ParseDeviceMapping(s)
}

func DeviceMappingFromStruct(name string, v cue.Value) (*api.DeviceMapping, error) {
	out := &api.DeviceMapping{Target: name}

	if p, ok := cuetils.Get(v, "target"); ok {
		if v, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.Target = v
		}
	}

	if p, ok := cuetils.Get(v, "source"); ok {
		if v, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.Source = v
		}
	}

	if p, ok := cuetils.Get(v, "permissions"); ok {
		if v, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.Permissions = v
		}
	}

	if p, ok := cuetils.Get(v, "cgroup_rule"); ok {
		if v, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.CgroupRule = v
		}
	}

	return out, nil
}
