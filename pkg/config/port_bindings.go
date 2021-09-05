package config

import (
	"strconv"

	"cuelang.org/go/cue"
	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/config"
)

func PortBindingsFromValue(v cue.Value) ([]*api.PortBinding, error) {
	if out, err := PortBindingsFromMap(v); err == nil {
		return out, nil
	}

	if out, err := PortBindingsFromList(v); err == nil {
		return out, nil
	}

	return nil, ErrUnexpectedSpec
}

func PortBindingsFromMap(v cue.Value) ([]*api.PortBinding, error) {
	out := []*api.PortBinding{}

	err := eachInMap(v, func(name string, v cue.Value) error {
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

	err := eachInList(v, func(v cue.Value) error {
		r, err := PortBindingFromValue("", v)
		if err == nil {
			out = append(out, r)
		}

		return err
	})

	return out, err
}

func PortBindingFromValue(name string, v cue.Value) (*api.PortBinding, error) {
	if out, err := PortBindingFromString(name, v); err == nil {
		return out, nil
	}

	if out, err := PortBindingFromStruct(name, v); err == nil {
		return out, nil
	}

	return nil, ErrUnexpectedSpec
}

func PortBindingFromString(_ string, v cue.Value) (*api.PortBinding, error) {
	s, err := v.String()
	if err != nil {
		return nil, err
	}

	return config.ParsePortBinding(s)
}

func PortBindingFromStruct(name string, v cue.Value) (*api.PortBinding, error) {
	out := &api.PortBinding{Target: name}

	if p, ok := property(v, "target"); ok {
		if v, err := p.String(); err == nil {
			out.Target = v
		} else if v, err := p.Int64(); err == nil {
			out.Target = strconv.FormatInt(v, 10)
		}
	}

	if p, ok := property(v, "publish"); ok {
		if v, err := p.String(); err == nil {
			out.Published = v
		} else if v, err := p.Int64(); err == nil {
			out.Published = strconv.FormatInt(v, 10)
		}
	}

	if p, ok := property(v, "protocol"); ok {
		if v, err := p.String(); err != nil {
			return nil, err
		} else {
			out.Protocol = v
		}
	}

	if p, ok := property(v, "host_ip"); ok {
		if v, err := p.String(); err != nil {
			return nil, err
		} else {
			out.HostIp = v
		}
	}

	return out, nil
}
