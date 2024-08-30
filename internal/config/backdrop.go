package config

import (
	"fmt"
	"strconv"

	"cuelang.org/go/cue"
	"github.com/wabenet/dodo-config/pkg/cuetils"
	api "github.com/wabenet/dodo-core/api/core/v1alpha5"
	"github.com/wabenet/dodo-core/pkg/config"
)

func BackdropFromStruct(name string, value cue.Value) (*api.Backdrop, error) {
	out := &api.Backdrop{
		Name:       name,
		Entrypoint: &api.Entrypoint{},
	}

	if p, ok, err := cuetils.Extract(value, "name", cuetils.String); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "name", err)
	} else if ok {
		out.Name = p
	}

	if p, ok, err := cuetils.Extract(value, "container_name", cuetils.String); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "container_name", err)
	} else if ok {
		out.ContainerName = p
	}

	if p, ok, err := cuetils.Extract(value, "runtime", cuetils.String); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "runtime", err)
	} else if ok {
		out.Runtime = p
	}

	if p, ok, err := cuetils.Extract(value, "script", cuetils.String); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "script", err)
	} else if ok {
		out.Entrypoint.Script = p
	}

	if p, ok, err := cuetils.Extract(value, "user", cuetils.String); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "uesr", err)
	} else if ok {
		out.User = p
	}

	if p, ok, err := cuetils.Extract(value, "working_dir", cuetils.String); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "working_dir", err)
	} else if ok {
		out.WorkingDir = p
	}

	if p, ok, err := cuetils.Extract(value, "image", cuetils.String); err != nil {
		if p, ok, err := cuetils.Extract(value, "image", BuildInfoFromStruct); err != nil {
			return nil, fmt.Errorf("invalid config for %s: %w", "image", err)
		} else if ok {
			out.BuildInfo = p
		}
	} else if ok {
		out.ImageId = p
	}

	if p, ok, err := cuetils.Extract(value, "build", cuetils.String); err != nil {
		if p, ok, err := cuetils.Extract(value, "build", BuildInfoFromStruct); err != nil {
			return nil, fmt.Errorf("invalid config for %s: %w", "build", err)
		} else if ok {
			out.BuildInfo = p
		}
	} else if ok {
		out.ImageId = p
	}

	if p, ok, err := cuetils.Extract(value, "aliases", cuetils.OneOrMore(cuetils.String)); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "aliases", err)
	} else if ok {
		out.Aliases = p
	}

	if p, ok, err := cuetils.Extract(value, "capabilities", cuetils.OneOrMore(cuetils.String)); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "capabilities", err)
	} else if ok {
		out.Capabilities = p
	}

	if p, ok, err := cuetils.Extract(value, "interpreter", cuetils.OneOrMore(cuetils.String)); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "interepreter", err)
	} else if ok {
		out.Entrypoint.Interpreter = p
	}

	if p, ok, err := cuetils.Extract(value, "environment", cuetils.ListOrDict(
		cuetils.Either([]cuetils.Extractor[*api.EnvironmentVariable]{
			cuetils.ParseString(config.ParseEnvironmentVariable),
			EnvironmentVariableFromStruct,
		}),
	)); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "environment", err)
	} else if ok {
		out.Environment = p
	}

	if p, ok, err := cuetils.Extract(value, "volumes", cuetils.ListOrDict(
		cuetils.Either([]cuetils.Extractor[*api.VolumeMount]{
			cuetils.ParseString(config.ParseVolumeMount),
			VolumeMountFromStruct,
		}),
	)); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "volumes", err)
	} else if ok {
		out.Volumes = p
	}

	if p, ok, err := cuetils.Extract(value, "ports", cuetils.ListOrDict(
		cuetils.Either([]cuetils.Extractor[*api.PortBinding]{
			cuetils.ParseString(config.ParsePortBinding),
			PortBindingFromStruct,
		}),
	)); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "ports", err)
	} else if ok {
		out.Ports = p
	}

	if p, ok, err := cuetils.Extract(value, "devices", cuetils.ListOrDict(
		cuetils.Either([]cuetils.Extractor[*api.DeviceMapping]{
			cuetils.ParseString(config.ParseDeviceMapping),
			DeviceMappingFromStruct,
		}),
	)); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "devices", err)
	} else if ok {
		out.Devices = p
	}

	return out, nil
}

func EnvironmentVariableFromStruct(name string, value cue.Value) (*api.EnvironmentVariable, error) {
	out := &api.EnvironmentVariable{Key: name}

	if p, ok, err := cuetils.Extract(value, "name", cuetils.String); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "name", err)
	} else if ok {
		out.Key = p
	}

	if p, ok, err := cuetils.Extract(value, "value", cuetils.String); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "value", err)
	} else if ok {
		out.Value = p
	}

	return out, nil
}

func VolumeMountFromStruct(name string, value cue.Value) (*api.VolumeMount, error) {
	out := &api.VolumeMount{Source: name}

	if p, ok, err := cuetils.Extract(value, "source", cuetils.String); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "source", err)
	} else if ok {
		out.Source = p
	}

	if p, ok, err := cuetils.Extract(value, "target", cuetils.String); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "target", err)
	} else if ok {
		out.Target = p
	}

	if p, ok, err := cuetils.Extract(value, "readonly", cuetils.Bool); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "readonly", err)
	} else if ok {
		out.Readonly = p
	}

	return out, nil
}

func PortBindingFromStruct(name string, value cue.Value) (*api.PortBinding, error) {
	out := &api.PortBinding{Target: name}

	if p, ok, err := cuetils.Extract(value, "target", cuetils.String); err != nil {
		if p, ok, err := cuetils.Extract(value, "target", cuetils.Int); err != nil {
			return nil, fmt.Errorf("invalid config for %s: %w", "target", err)
		} else if ok {
			out.Target = strconv.FormatInt(p, 10)
		}
	} else if ok {
		out.Target = p
	}

	if p, ok, err := cuetils.Extract(value, "publish", cuetils.String); err != nil {
		if p, ok, err := cuetils.Extract(value, "publish", cuetils.Int); err != nil {
			return nil, fmt.Errorf("invalid config for %s: %w", "publish", err)
		} else if ok {
			out.Published = strconv.FormatInt(p, 10)
		}
	} else if ok {
		out.Published = p
	}

	if p, ok, err := cuetils.Extract(value, "protocol", cuetils.String); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "protocol", err)
	} else if ok {
		out.Protocol = p
	}

	if p, ok, err := cuetils.Extract(value, "host_ip", cuetils.String); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "host_ip", err)
	} else if ok {
		out.HostIp = p
	}

	return out, nil
}

func DeviceMappingFromStruct(name string, value cue.Value) (*api.DeviceMapping, error) {
	out := &api.DeviceMapping{Target: name}

	if p, ok, err := cuetils.Extract(value, "target", cuetils.String); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "target", err)
	} else if ok {
		out.Target = p
	}

	if p, ok, err := cuetils.Extract(value, "source", cuetils.String); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "source", err)
	} else if ok {
		out.Source = p
	}

	if p, ok, err := cuetils.Extract(value, "permissions", cuetils.String); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "permissions", err)
	} else if ok {
		out.Permissions = p
	}

	if p, ok, err := cuetils.Extract(value, "cgroup_rule", cuetils.String); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "cgroup_rule", err)
	} else if ok {
		out.CgroupRule = p
	}

	return out, nil
}
