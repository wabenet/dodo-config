package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wabenet/dodo-config/internal/config"
	"github.com/wabenet/dodo-core/pkg/plugin/configuration"
	"github.com/wabenet/dodo-core/pkg/plugin/runtime"
)

func TestAllDefaults(t *testing.T) {
	loadBackdrop(t, "test_all_defaults")
}

func TestMinus(t *testing.T) {
	loadBackdrop(t, "test-minus")
}

func TestBasicBackdrop(t *testing.T) {
	backdrop := loadBackdrop(t, "test_full_configs")

	assert.Equal(t, "testimage", backdrop.ContainerConfig.Image)
	assert.Equal(t, "testcontainer", backdrop.ContainerConfig.Name)
	assert.Equal(t, "/home/test", backdrop.ContainerConfig.Process.WorkingDir)
	assert.Equal(t, runtime.Entrypoint{"/bin/sh", backdrop.RequiredFiles[1].FilePath}, backdrop.ContainerConfig.Process.Entrypoint)
	assert.Equal(t, []byte("echo \"$@\"\n"), backdrop.RequiredFiles[1].Contents)
}

func TestBuildInfo(t *testing.T) {
	backdrop := loadBackdrop(t, "test_build")

	build := backdrop.BuildConfig
	assert.NotNil(t, backdrop.BuildConfig)

	assert.Equal(t, "testimage", build.ImageName)
	assert.Equal(t, "/some/path", build.Context)
	assert.Equal(t, "/some/other/path", build.Dockerfile)
	assert.Equal(t, []string{"FROM foo\n"}, build.InlineDockerfile)
}

func TestFullEnvironment(t *testing.T) {
	backdrop := loadBackdrop(t, "test_full_configs")

	for _, env := range backdrop.ContainerConfig.Environment {
		assert.NotEqual(t, "FULL", env.Key)
	}

	assert.Contains(t, backdrop.ContainerConfig.Environment, runtime.EnvironmentVariable{
		Key:   "FOO",
		Value: "BAR",
	})
}

func TestPartialEnvironment(t *testing.T) {
	backdrop := loadBackdrop(t, "test_full_configs")

	assert.Contains(t, backdrop.ContainerConfig.Environment, runtime.EnvironmentVariable{
		Key:   "PARTIAL",
		Value: "",
	})
}

func TestEnvironmentWithList(t *testing.T) {
	backdrop := loadBackdrop(t, "test_with_lists")

	assert.Contains(t, backdrop.ContainerConfig.Environment, runtime.EnvironmentVariable{
		Key:   "FOO",
		Value: "BAR",
	})

	assert.Contains(t, backdrop.ContainerConfig.Environment, runtime.EnvironmentVariable{
		Key:   "SOMETHING",
		Value: "",
	})
}

func TestFullPortBindings(t *testing.T) {
	backdrop := loadBackdrop(t, "test_full_configs")

	assert.Contains(t, backdrop.ContainerConfig.Ports, runtime.PortBinding{
		ContainerPort: "80",
		HostPort:      "8080",
		HostIP:        "192.168.0.1",
	})
}

func TestPortBindingsWithList(t *testing.T) {
	backdrop := loadBackdrop(t, "test_with_lists")

	assert.Contains(t, backdrop.ContainerConfig.Ports, runtime.PortBinding{
		ContainerPort: "80",
		HostPort:      "8080",
	})
}

func TestFullFiles(t *testing.T) {
	backdrop := loadBackdrop(t, "test_full_configs")

	assert.Contains(t, backdrop.RequiredFiles, configuration.File{
		FilePath: "/foo/hello.txt",
		Contents: []byte("Hello World!\n"),
	})
}

func TestFilesWithLists(t *testing.T) {
	backdrop := loadBackdrop(t, "test_with_lists")

	assert.Contains(t, backdrop.RequiredFiles, configuration.File{
		FilePath: "/foo/hello.txt",
		Contents: []byte("Hello World!"),
	})
}

func TestFullMounts(t *testing.T) {
	backdrop := loadBackdrop(t, "test_full_configs")

	assert.Contains(t, backdrop.ContainerConfig.Mounts, runtime.BindMount{
		HostPath:      "/from/path",
		ContainerPath: "/to/path",
		Readonly:      true,
	})

	assert.Contains(t, backdrop.ContainerConfig.Mounts, runtime.BindMount{
		HostPath:      "/some/mount",
		ContainerPath: "",
		Readonly:      false,
	})

	assert.Contains(t, backdrop.ContainerConfig.Mounts, runtime.DeviceMount{
		HostPath:      "/dev/snd",
		ContainerPath: "/foo/bar",
		Permissions:   "rw",
	})

	assert.Contains(t, backdrop.ContainerConfig.Mounts, runtime.DeviceMount{
		ContainerPath: "rule", // TODO why?
		CGroupRule:    "c *:* rmw",
	})
}

func TestMountsWithLists(t *testing.T) {
	backdrop := loadBackdrop(t, "test_with_lists")

	assert.Contains(t, backdrop.ContainerConfig.Mounts, runtime.VolumeMount{
		VolumeName:    "foo",
		ContainerPath: "bar",
		Readonly:      true,
	})

	assert.Contains(t, backdrop.ContainerConfig.Mounts, runtime.DeviceMount{
		HostPath:      "/dev/snd",
		ContainerPath: "/dev/snd",
	})
}

func loadBackdrop(t *testing.T, name string) configuration.Backdrop {
	cfg, err := config.GetAllBackdrops("test/dodo.yaml")
	require.NoError(t, err)

	backdrop, ok := cfg[name]
	assert.True(t, ok)

	return backdrop
}
