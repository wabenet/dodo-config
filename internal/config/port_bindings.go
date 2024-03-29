package config

import (
	"strconv"

	"cuelang.org/go/cue"
	"github.com/hashicorp/go-multierror"
	"github.com/wabenet/dodo-config/pkg/cuetils"
	api "github.com/wabenet/dodo-core/api/core/v1alpha5"
	"github.com/wabenet/dodo-core/pkg/config"
)

func PortBindingsFromValue(v cue.Value) ([]*api.PortBinding, error) {
	var errs error

	if out, err := PortBindingsFromMap(v); err == nil {
		return out, nil
	} else {
		errs = multierror.Append(errs, err)
	}

	if out, err := PortBindingsFromList(v); err == nil {
		return out, nil
	} else {
		errs = multierror.Append(errs, err)
	}

	return nil, errs
}

func PortBindingsFromMap(v cue.Value) ([]*api.PortBinding, error) {
	out := []*api.PortBinding{}

	err := cuetils.IterMap(v, func(name string, v cue.Value) error {
		r, err := PortBindingFromValue(name, v)
		if err == nil {
			out = append(out, r)
		}

		return err

	})

	return out, err
}

func PortBindingsFromList(v cue.Value) ([]*api.PortBinding, error) {
	out := []*api.PortBinding{}

	err := cuetils.IterList(v, func(v cue.Value) error {
		r, err := PortBindingFromValue("", v)
		if err == nil {
			out = append(out, r)
		}

		return err
	})

	return out, err
}

func PortBindingFromValue(name string, v cue.Value) (*api.PortBinding, error) {
	var errs error

	if out, err := PortBindingFromString(name, v); err == nil {
		return out, nil
	} else {
		errs = multierror.Append(errs, err)
	}

	if out, err := PortBindingFromStruct(name, v); err == nil {
		return out, nil
	} else {
		errs = multierror.Append(errs, err)
	}

	return nil, errs
}

func PortBindingFromString(_ string, v cue.Value) (*api.PortBinding, error) {
	s, err := StringFromValue(v)
	if err != nil {
		return nil, err
	}

	return config.ParsePortBinding(s)
}

func PortBindingFromStruct(name string, v cue.Value) (*api.PortBinding, error) {
	out := &api.PortBinding{Target: name}

	if p, ok := cuetils.Get(v, "target"); ok {
		if v, err := StringFromValue(p); err == nil {
			out.Target = v
		} else if v, err := p.Int64(); err == nil {
			out.Target = strconv.FormatInt(v, 10)
		}
	}

	if p, ok := cuetils.Get(v, "publish"); ok {
		if v, err := StringFromValue(p); err == nil {
			out.Published = v
		} else if v, err := p.Int64(); err == nil {
			out.Published = strconv.FormatInt(v, 10)
		}
	}

	if p, ok := cuetils.Get(v, "protocol"); ok {
		if v, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.Protocol = v
		}
	}

	if p, ok := cuetils.Get(v, "host_ip"); ok {
		if v, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.HostIp = v
		}
	}

	return out, nil
}
