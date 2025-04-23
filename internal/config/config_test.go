package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wabenet/dodo-config/internal/config"
	api "github.com/wabenet/dodo-core/api/configuration/v1alpha2"
	runtime "github.com/wabenet/dodo-core/api/runtime/v1alpha2"
	"google.golang.org/protobuf/proto"
)

func TestAllDefaults(t *testing.T) {
	loadBackdrop(t, "test_all_defaults")
}

func TestMinus(t *testing.T) {
	loadBackdrop(t, "test-minus")
}

func TestBasicBackdrop(t *testing.T) {
	backdrop := loadBackdrop(t, "test_full_configs")

	assert.Equal(t, "testimage", backdrop.GetContainerConfig().GetImage())
	assert.Equal(t, "testcontainer", backdrop.GetContainerConfig().GetName())
	assert.Equal(t, "/home/test", backdrop.GetContainerConfig().GetProcess().GetWorkingDir())
	assert.Equal(t, []string{"/bin/sh", backdrop.GetRequiredFiles()[0].GetFilePath()}, backdrop.GetContainerConfig().GetProcess().GetEntrypoint())
	assert.Equal(t, "echo \"$@\"\n", backdrop.GetRequiredFiles()[0].GetContents())
}

func TestBuildInfo(t *testing.T) {
	backdrop := loadBackdrop(t, "test_build")

	build := backdrop.GetBuildConfig()
	assert.NotNil(t, backdrop.GetBuildConfig())

	assert.Equal(t, "testimage", build.GetImageName())
	assert.Equal(t, "/some/path", build.GetContext())
	assert.Equal(t, "/some/other/path", build.GetDockerfile())
	assert.Equal(t, []string{"FROM foo\n"}, build.GetInlineDockerfile())
}

func TestFullEnvironment(t *testing.T) {
	backdrop := loadBackdrop(t, "test_full_configs")

	for _, env := range backdrop.GetContainerConfig().GetEnvironment() {
		assert.NotEqual(t, "FULL", env.GetKey())
	}

	assertContainsEnvironment(t, backdrop, &runtime.EnvironmentVariable{
		Key:   "FOO",
		Value: "BAR",
	})
}

func TestPartialEnvironment(t *testing.T) {
	backdrop := loadBackdrop(t, "test_full_configs")

	assertContainsEnvironment(t, backdrop, &runtime.EnvironmentVariable{
		Key:   "PARTIAL",
		Value: "",
	})
}

func TestEnvironmentWithList(t *testing.T) {
	backdrop := loadBackdrop(t, "test_with_lists")

	assertContainsEnvironment(t, backdrop, &runtime.EnvironmentVariable{
		Key:   "FOO",
		Value: "BAR",
	})

	assertContainsEnvironment(t, backdrop, &runtime.EnvironmentVariable{
		Key:   "SOMETHING",
		Value: "",
	})
}

func TestFullPortBindings(t *testing.T) {
	backdrop := loadBackdrop(t, "test_full_configs")

	assertContainsPortBinding(t, backdrop, &runtime.PortBinding{
		ContainerPort: "80",
		HostPort:      "8080",
		HostIp:        "192.168.0.1",
	})
}

func TestPortBindingsWithList(t *testing.T) {
	backdrop := loadBackdrop(t, "test_with_lists")

	assertContainsPortBinding(t, backdrop, &runtime.PortBinding{
		ContainerPort: "80",
		HostPort:      "8080",
	})
}

func TestFullMounts(t *testing.T) {
	backdrop := loadBackdrop(t, "test_full_configs")

	assertContainsMount(t, backdrop, &runtime.Mount{
		Type: &runtime.Mount_Bind{
			Bind: &runtime.BindMount{
				HostPath:      "/from/path",
				ContainerPath: "/to/path",
				Readonly:      true,
			},
		},
	})

	assertContainsMount(t, backdrop, &runtime.Mount{
		Type: &runtime.Mount_Bind{
			Bind: &runtime.BindMount{
				HostPath:      "/some/mount",
				ContainerPath: "",
				Readonly:      false,
			},
		},
	})

	assertContainsMount(t, backdrop, &runtime.Mount{
		Type: &runtime.Mount_Device{
			Device: &runtime.DeviceMount{
				HostPath:      "/dev/snd",
				ContainerPath: "/foo/bar",
				Permissions:   "rw",
			},
		},
	})

	assertContainsMount(t, backdrop, &runtime.Mount{
		Type: &runtime.Mount_Device{
			Device: &runtime.DeviceMount{
				ContainerPath: "rule", // TODO why?
				CgroupRule:    "c *:* rmw",
			},
		},
	})
}

func TestMountsWithLists(t *testing.T) {
	backdrop := loadBackdrop(t, "test_with_lists")

	assertContainsMount(t, backdrop, &runtime.Mount{
		Type: &runtime.Mount_Volume{
			Volume: &runtime.VolumeMount{
				VolumeName:    "foo",
				ContainerPath: "bar",
				Readonly:      true,
			},
		},
	})

	assertContainsMount(t, backdrop, &runtime.Mount{
		Type: &runtime.Mount_Device{
			Device: &runtime.DeviceMount{
				HostPath:      "/dev/snd",
				ContainerPath: "/dev/snd",
			},
		},
	})
}

func loadBackdrop(t *testing.T, name string) *api.Backdrop {
	cfg, err := config.GetAllBackdrops("test/dodo.yaml")
	require.NoError(t, err)

	backdrop, ok := cfg[name]
	assert.True(t, ok)

	return backdrop
}

func assertContainsEnvironment(t *testing.T, backdrop *api.Backdrop, expected *runtime.EnvironmentVariable) bool {
	assert.NotEmpty(t, backdrop.GetContainerConfig().GetEnvironment())

	for _, env := range backdrop.GetContainerConfig().GetEnvironment() {
		if proto.Equal(env, expected) {
			return true
		}
	}

	return assert.Fail(t, "did not find expected environment config")
}

func assertContainsPortBinding(t *testing.T, backdrop *api.Backdrop, expected *runtime.PortBinding) bool {
	assert.NotEmpty(t, backdrop.GetContainerConfig().GetPorts())

	for _, port := range backdrop.GetContainerConfig().GetPorts() {
		if proto.Equal(port, expected) {
			return true
		}
	}

	return assert.Fail(t, "did not find expected port config")
}

func assertContainsMount(t *testing.T, backdrop *api.Backdrop, expected *runtime.Mount) bool {
	assert.NotEmpty(t, backdrop.GetContainerConfig().GetMounts())

	for _, mnt := range backdrop.GetContainerConfig().GetMounts() {
		if proto.Equal(mnt, expected) {
			return true
		}
	}

	return assert.Fail(t, "did not find expected mount config")
}
