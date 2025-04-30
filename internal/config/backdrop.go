package config

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"cuelang.org/go/cue"
	"github.com/docker/docker/pkg/stringid"
	"github.com/wabenet/dodo-config/pkg/cuetils"
	"github.com/wabenet/dodo-core/pkg/plugin/builder"
	"github.com/wabenet/dodo-core/pkg/plugin/configuration"
	"github.com/wabenet/dodo-core/pkg/plugin/runtime"
)

func BackdropFromStruct(name string, value cue.Value) (configuration.Backdrop, error) {
	out := configuration.Backdrop{Name: name}

	if p, ok, err := cuetils.Extract(value, "name", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "name", err)
	} else if ok {
		out.Name = p
	}

	if p, ok, err := cuetils.Extract(value, "aliases", cuetils.OneOrMore(cuetils.String)); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "aliases", err)
	} else if ok {
		out.Aliases = p
	}

	if p, ok, err := cuetils.Extract(value, "runtime", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "runtime", err)
	} else if ok {
		out.Runtime = p
	}

	if p, ok, err := cuetils.Extract(value, "build.builder", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "builder", err)
	} else if ok {
		out.Builder = p
	}

	if p, err := ContainerConfigFromStruct(name, value); err != nil {
		return out, err
	} else {
		out.ContainerConfig = p
	}

	if p, ok, err := cuetils.Extract(value, "image", cuetils.String); err != nil {
		if p, ok, err := cuetils.Extract(value, "image", BuildConfigFromStruct); err != nil {
			return out, fmt.Errorf("invalid config for %s: %w", "image", err)
		} else if ok {
			out.BuildConfig = p
		}
	} else if ok {
		out.ContainerConfig.Image = p
	}

	if p, ok, err := cuetils.Extract(value, "build", cuetils.String); err != nil {
		if p, ok, err := cuetils.Extract(value, "build", BuildConfigFromStruct); err != nil {
			return out, fmt.Errorf("invalid config for %s: %w", "build", err)
		} else if ok {
			out.BuildConfig = p
		}
	} else if ok {
		out.ContainerConfig.Image = p
	}

	if p, ok, err := cuetils.Extract(value, "script", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "script", err)
	} else if ok {
		tmpPath := fmt.Sprintf("/tmp/dodo-%s/", stringid.GenerateRandomID()[:20])
		out.RequiredFiles = append(out.RequiredFiles, configuration.File{
			FilePath: path.Join(tmpPath, "entrypoint"),
			Contents: []byte(p),
		})
		out.ContainerConfig.Process.Entrypoint = append(out.ContainerConfig.Process.Entrypoint, path.Join(tmpPath, "entrypoint"))
	}

	return out, nil
}

func ContainerConfigFromStruct(name string, value cue.Value) (runtime.ContainerConfig, error) {
	out := runtime.ContainerConfig{
		Mounts: []runtime.Mount{},
	}

	if p, ok, err := cuetils.Extract(value, "container_name", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "container_name", err)
	} else if ok {
		out.Name = p
	}

	if p, err := ProcessFromStruct(name, value); err != nil {
		return out, err
	} else {
		out.Process = p
	}

	if p, ok, err := cuetils.Extract(value, "capabilities", cuetils.OneOrMore(cuetils.String)); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "capabilities", err)
	} else if ok {
		out.Capabilities = p
	}

	if p, ok, err := cuetils.Extract(value, "environment", cuetils.ListOrDict(
		cuetils.Either([]cuetils.Extractor[runtime.EnvironmentVariable]{
			cuetils.ParseString(runtime.EnvironmentVariableFromSpec),
			EnvironmentVariableFromStruct,
		}),
	)); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "environment", err)
	} else if ok {
		out.Environment = p
	}

	if p, ok, err := cuetils.Extract(value, "ports", cuetils.ListOrDict(
		cuetils.Either([]cuetils.Extractor[runtime.PortBinding]{
			cuetils.ParseString(runtime.PortBindingFromSpec),
			PortBindingFromStruct,
		}),
	)); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "ports", err)
	} else if ok {
		out.Ports = p
	}

	if p, ok, err := cuetils.Extract(value, "mounts", cuetils.ListOrDict(
		cuetils.Either([]cuetils.Extractor[runtime.Mount]{
			BindMountFromStruct,
			VolumeMountFromStruct,
			TmpfsMountFromStruct,
			ImageMountFromStruct,
			DeviceMountFromStruct,
		}),
	)); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "volumes", err)
	} else if ok {
		out.Mounts = append(out.Mounts, p...)
	}

	// Deprecated
	if p, ok, err := cuetils.Extract(value, "volumes", cuetils.ListOrDict(
		cuetils.Either([]cuetils.Extractor[runtime.Mount]{
			cuetils.ParseString(runtime.BindMountFromSpec),
			VolumeMountFromStruct,
			BindMountFromStruct,
		}),
	)); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "volumes", err)
	} else if ok {
		out.Mounts = append(out.Mounts, p...)
	}

	// Deprecated
	if p, ok, err := cuetils.Extract(value, "devices", cuetils.ListOrDict(
		cuetils.Either([]cuetils.Extractor[runtime.Mount]{
			cuetils.ParseString(runtime.DeviceMountFromSpec),
			DeviceMountFromStruct,
		}),
	)); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "devices", err)
	} else if ok {
		out.Mounts = append(out.Mounts, p...)
	}

	return out, nil
}

func ProcessFromStruct(name string, value cue.Value) (runtime.Process, error) {
	out := runtime.Process{}

	if p, ok, err := cuetils.Extract(value, "interpreter", cuetils.OneOrMore(cuetils.String)); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "interepreter", err)
	} else if ok {
		out.Entrypoint = p
	} else {
		out.Entrypoint = []string{"/bin/sh"}
	}

	if p, ok, err := cuetils.Extract(value, "user", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "uesr", err)
	} else if ok {
		out.User = p
	}

	if p, ok, err := cuetils.Extract(value, "working_dir", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "working_dir", err)
	} else if ok {
		out.WorkingDir = p
	}

	return out, nil
}

func EnvironmentVariableFromStruct(name string, value cue.Value) (runtime.EnvironmentVariable, error) {
	out := runtime.EnvironmentVariable{Key: name}

	if p, ok, err := cuetils.Extract(value, "name", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "name", err)
	} else if ok {
		out.Key = p
	}

	if p, ok, err := cuetils.Extract(value, "value", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "value", err)
	} else if ok {
		out.Value = p
	}

	return out, nil
}

func PortBindingFromStruct(name string, value cue.Value) (runtime.PortBinding, error) {
	out := runtime.PortBinding{HostPort: name}

	if p, ok, err := cuetils.Extract(value, "target", cuetils.String); err != nil {
		if p, ok, err := cuetils.Extract(value, "target", cuetils.Int); err != nil {
			return out, fmt.Errorf("invalid config for %s: %w", "target", err)
		} else if ok {
			out.ContainerPort = strconv.FormatInt(p, 10)
		}
	} else if ok {
		out.HostPort = p
	}

	if p, ok, err := cuetils.Extract(value, "publish", cuetils.String); err != nil {
		if p, ok, err := cuetils.Extract(value, "publish", cuetils.Int); err != nil {
			return out, fmt.Errorf("invalid config for %s: %w", "publish", err)
		} else if ok {
			out.HostPort = strconv.FormatInt(p, 10)
		}
	} else if ok {
		out.ContainerPort = p
	}

	if p, ok, err := cuetils.Extract(value, "protocol", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "protocol", err)
	} else if ok {
		out.Protocol = p
	}

	if p, ok, err := cuetils.Extract(value, "host_ip", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "host_ip", err)
	} else if ok {
		out.HostIP = p
	}

	return out, nil
}

func BindMountFromStruct(name string, value cue.Value) (runtime.Mount, error) {
	out := runtime.BindMount{HostPath: name}

	if p, ok, err := cuetils.Extract(value, "type", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "type", err)
	} else if !ok {
		return out, fmt.Errorf("%s is missing required type on value %v", name, value)
	} else if p != "bind" {
		return out, fmt.Errorf("%s is not a bind mount config, but %s", name, p)
	}

	if p, ok, err := cuetils.Extract(value, "source", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "source", err)
	} else if ok {
		out.HostPath = p
	}

	if p, ok, err := cuetils.Extract(value, "target", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "target", err)
	} else if ok {
		out.ContainerPath = p
	}

	if p, ok, err := cuetils.Extract(value, "readonly", cuetils.Bool); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "readonly", err)
	} else if ok {
		out.Readonly = p
	}

	return out, nil
}

func VolumeMountFromStruct(name string, value cue.Value) (runtime.Mount, error) {
	out := runtime.VolumeMount{VolumeName: name}

	if p, ok, err := cuetils.Extract(value, "type", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "type", err)
	} else if !ok || p != "volume" {
		return out, fmt.Errorf("%s is not a volume mount config", name)
	}

	if p, ok, err := cuetils.Extract(value, "source", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "source", err)
	} else if ok {
		out.VolumeName = p
	}

	if p, ok, err := cuetils.Extract(value, "target", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "target", err)
	} else if ok {
		out.ContainerPath = p
	}

	if p, ok, err := cuetils.Extract(value, "path", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "path", err)
	} else if ok {
		out.Subpath = p
	}

	if p, ok, err := cuetils.Extract(value, "readonly", cuetils.Bool); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "readonly", err)
	} else if ok {
		out.Readonly = p
	}

	return out, nil
}

func TmpfsMountFromStruct(name string, value cue.Value) (runtime.Mount, error) {
	out := runtime.TmpfsMount{ContainerPath: name}

	if p, ok, err := cuetils.Extract(value, "type", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "type", err)
	} else if !ok || p != "tmpfs" {
		return out, fmt.Errorf("%s is not a tmpfs mount config", name)
	}

	if p, ok, err := cuetils.Extract(value, "target", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "target", err)
	} else if ok {
		out.ContainerPath = p
	}

	if p, ok, err := cuetils.Extract(value, "size", cuetils.Int); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "path", err)
	} else if ok {
		out.Size = int(p)
	}

	if p, ok, err := cuetils.Extract(value, "mode", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "readonly", err)
	} else if ok {
		mode, mErr := strconv.ParseUint(p, 8, 32)
		if mErr != nil {
			return out, fmt.Errorf("invalid file mode for %s: %w", name, mErr)
		}

		out.Mode = os.FileMode(mode)
	}

	return out, nil
}

func ImageMountFromStruct(name string, value cue.Value) (runtime.Mount, error) {
	out := runtime.ImageMount{Image: name}

	if p, ok, err := cuetils.Extract(value, "type", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "type", err)
	} else if !ok || p != "image" {
		return out, fmt.Errorf("%s is not an image mount config", name)
	}

	if p, ok, err := cuetils.Extract(value, "source", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "source", err)
	} else if ok {
		out.Image = p
	}

	if p, ok, err := cuetils.Extract(value, "target", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "target", err)
	} else if ok {
		out.ContainerPath = p
	}

	if p, ok, err := cuetils.Extract(value, "path", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "path", err)
	} else if ok {
		out.Subpath = p
	}

	if p, ok, err := cuetils.Extract(value, "readonly", cuetils.Bool); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "readonly", err)
	} else if ok {
		out.Readonly = p
	}

	return out, nil
}

func DeviceMountFromStruct(name string, value cue.Value) (runtime.Mount, error) {
	out := runtime.DeviceMount{ContainerPath: name}

	if p, ok, err := cuetils.Extract(value, "type", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "type", err)
	} else if !ok || p != "device" {
		return out, fmt.Errorf("%s is not a device mount config", name)
	}

	if p, ok, err := cuetils.Extract(value, "target", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "target", err)
	} else if ok {
		out.ContainerPath = p
	}

	if p, ok, err := cuetils.Extract(value, "source", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "source", err)
	} else if ok {
		out.HostPath = p
	}

	if p, ok, err := cuetils.Extract(value, "permissions", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "permissions", err)
	} else if ok {
		out.Permissions = p
	}

	if p, ok, err := cuetils.Extract(value, "cgroup_rule", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "cgroup_rule", err)
	} else if ok {
		out.CGroupRule = p
	}

	return out, nil
}

func BuildConfigFromStruct(_ string, value cue.Value) (builder.BuildConfig, error) {
	out := builder.BuildConfig{}

	if p, ok, err := cuetils.Extract(value, "name", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "name", err)
	} else if ok {
		out.ImageName = p
	}

	if p, ok, err := cuetils.Extract(value, "context", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "context", err)
	} else if ok {
		out.Context = p
	}

	if p, ok, err := cuetils.Extract(value, "dockerfile", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "dockerfile", err)
	} else if ok {
		out.Dockerfile = p
	}

	if p, ok, err := cuetils.Extract(value, "steps", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "steps", err)
	} else if ok {
		out.InlineDockerfile = []string{p}
	}

	if p, ok, err := cuetils.Extract(value, "dependencies", cuetils.OneOrMore(cuetils.String)); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "dependencies", err)
	} else if ok {
		out.Dependencies = p
	}

	if p, ok, err := cuetils.Extract(value, "arguments", cuetils.ListOrDict(BuildArgumentFromStruct)); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "arguments", err)
	} else if ok {
		out.Arguments = p
	}

	if p, ok, err := cuetils.Extract(value, "secrets", cuetils.ListOrDict(BuildSecretFromStruct)); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "secrets", err)
	} else if ok {
		out.Secrets = p
	}

	if p, ok, err := cuetils.Extract(value, "ssh_agents", cuetils.ListOrDict(BuildSSHAgentFromStruct)); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "ssh_agents", err)
	} else if ok {
		out.SSHAgents = p
	}

	return out, nil
}

func BuildArgumentFromStruct(name string, value cue.Value) (builder.BuildArgument, error) {
	out := builder.BuildArgument{Key: name}

	if p, ok, err := cuetils.Extract(value, "name", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "name", err)
	} else if ok {
		out.Key = p
	}

	if p, ok, err := cuetils.Extract(value, "value", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "value", err)
	} else if ok {
		out.Value = p
	}

	return out, nil
}

func BuildSecretFromStruct(name string, value cue.Value) (builder.BuildSecret, error) {
	out := builder.BuildSecret{ID: name}

	if p, ok, err := cuetils.Extract(value, "id", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "id", err)
	} else if ok {
		out.ID = p
	}

	if p, ok, err := cuetils.Extract(value, "path", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "path", err)
	} else if ok {
		out.Path = p
	}

	return out, nil
}

func BuildSSHAgentFromStruct(name string, value cue.Value) (builder.SSHAgent, error) {
	out := builder.SSHAgent{ID: name}

	if p, ok, err := cuetils.Extract(value, "path", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "path", err)
	} else if ok {
		out.ID = p
	}

	if p, ok, err := cuetils.Extract(value, "identity_file", cuetils.String); err != nil {
		return out, fmt.Errorf("invalid config for %s: %w", "identity_file", err)
	} else if ok {
		out.IdentityFile = p
	}

	return out, nil
}
