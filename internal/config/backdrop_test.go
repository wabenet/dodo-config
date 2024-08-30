package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wabenet/dodo-config/internal/config"
)

func TestAllDefaults(t *testing.T) {
	cfg, err := config.GetAllBackdrops("test/dodo.yaml")
	require.NoError(t, err)

	_, ok := cfg["test_all_defaults"]
	assert.True(t, ok)
}

func TestMinus(t *testing.T) {
	cfg, err := config.GetAllBackdrops("test/dodo.yaml")
	require.NoError(t, err)

	_, ok := cfg["test-minus"]
	assert.True(t, ok)
}

func TestBasicBackdrop(t *testing.T) {
	cfg, err := config.GetAllBackdrops("test/dodo.yaml")
	require.NoError(t, err)

	backdrop, ok := cfg["test_full_configs"]

	assert.True(t, ok)
	assert.Equal(t, "testimage", backdrop.GetImageId())
	assert.Equal(t, "testcontainer", backdrop.GetContainerName())
	assert.Equal(t, "/home/test", backdrop.GetWorkingDir())
	assert.Equal(t, "echo \"$@\"\n", backdrop.GetEntrypoint().GetScript())
}

func TestFullEnvironment(t *testing.T) {
	cfg, err := config.GetAllBackdrops("test/dodo.yaml")
	require.NoError(t, err)

	backdrop, ok := cfg["test_full_configs"]

	assert.True(t, ok)
	assert.NotEmpty(t, backdrop.GetEnvironment())

	for _, env := range backdrop.GetEnvironment() {
		assert.NotEqual(t, "FULL", env.GetKey())

		if env.GetKey() == "FOO" {
			assert.Equal(t, "BAR", env.GetValue())

			return
		}
	}

	assert.Fail(t, "did not find expected environment config 'FULL'")
}

func TestPartialEnvironment(t *testing.T) {
	cfg, err := config.GetAllBackdrops("test/dodo.yaml")
	require.NoError(t, err)

	backdrop, ok := cfg["test_full_configs"]

	assert.True(t, ok)
	assert.NotEmpty(t, backdrop.GetEnvironment())

	for _, env := range backdrop.GetEnvironment() {
		if env.GetKey() == "PARTIAL" {
			assert.Equal(t, "", env.GetValue())

			return
		}
	}

	assert.Fail(t, "did not find expected environment config 'PARTIAL'")
}

func TestEnvironmentWithList(t *testing.T) {
	cfg, err := config.GetAllBackdrops("test/dodo.yaml")
	require.NoError(t, err)

	backdrop, ok := cfg["test_with_lists"]

	assert.True(t, ok)
	assert.Len(t, backdrop.GetEnvironment(), 2)

	for _, env := range backdrop.GetEnvironment() {
		switch env.GetKey() {
		case "FOO":
			assert.Equal(t, "BAR", env.GetValue())
		case "SOMETHING":
			assert.Equal(t, "", env.GetValue())
		default:
			assert.Fail(t, "unexpected environment config")
		}
	}
}

func TestFullVolume(t *testing.T) {
	cfg, err := config.GetAllBackdrops("test/dodo.yaml")
	require.NoError(t, err)

	backdrop, ok := cfg["test_full_configs"]

	assert.True(t, ok)
	assert.NotEmpty(t, backdrop.GetVolumes())

	for _, vol := range backdrop.GetVolumes() {
		if vol.GetSource() == "/from/path" {
			assert.Equal(t, "/to/path", vol.GetTarget())
			assert.True(t, vol.GetReadonly())

			return
		}
	}

	assert.Fail(t, "did not find expected volume config 'full'")
}

func TestPartialVolume(t *testing.T) {
	cfg, err := config.GetAllBackdrops("test/dodo.yaml")
	require.NoError(t, err)

	backdrop, ok := cfg["test_full_configs"]

	assert.True(t, ok)
	assert.NotEmpty(t, backdrop.GetVolumes())

	for _, vol := range backdrop.GetVolumes() {
		if vol.GetSource() == "/some/mount" {
			return
		}
	}

	assert.Fail(t, "did not find expected volume config 'full'")
}

func TestVolumesWithLists(t *testing.T) {
	cfg, err := config.GetAllBackdrops("test/dodo.yaml")
	require.NoError(t, err)

	backdrop, ok := cfg["test_with_lists"]

	assert.True(t, ok)
	assert.Len(t, backdrop.GetVolumes(), 1)

	for _, vol := range backdrop.GetVolumes() {
		switch vol.GetSource() {
		case "foo":
			assert.Equal(t, "bar", vol.GetTarget())
			assert.True(t, vol.GetReadonly())
		default:
			assert.Fail(t, "unexpected volume config")
		}
	}
}

func TestFullPortBindings(t *testing.T) {
	cfg, err := config.GetAllBackdrops("test/dodo.yaml")
	require.NoError(t, err)

	backdrop, ok := cfg["test_full_configs"]

	assert.True(t, ok)
	assert.Len(t, backdrop.GetPorts(), 1)

	for _, port := range backdrop.GetPorts() {
		switch port.GetTarget() {
		case "80":
			assert.Equal(t, "8080", port.GetPublished())
			assert.Equal(t, "192.168.0.1", port.GetHostIp())
		default:
			assert.Fail(t, "unexpected port config")
		}
	}
}

func TestPortBindingsWithList(t *testing.T) {
	cfg, err := config.GetAllBackdrops("test/dodo.yaml")
	require.NoError(t, err)

	backdrop, ok := cfg["test_with_lists"]

	assert.True(t, ok)
	assert.Len(t, backdrop.GetPorts(), 1)

	for _, port := range backdrop.GetPorts() {
		switch port.GetTarget() {
		case "80":
			assert.Equal(t, "8080", port.GetPublished())
		default:
			assert.Fail(t, "unexpected port config")
		}
	}
}

func TestFullDeviceMappings(t *testing.T) {
	cfg, err := config.GetAllBackdrops("test/dodo.yaml")
	require.NoError(t, err)

	backdrop, ok := cfg["test_full_configs"]

	assert.True(t, ok)
	assert.Len(t, backdrop.GetDevices(), 2)

	for _, dev := range backdrop.GetDevices() {
		switch dev.GetTarget() {
		case "/foo/bar":
			assert.Equal(t, "/dev/snd", dev.GetSource())
			assert.Equal(t, "rw", dev.GetPermissions())
		case "rule":
			assert.Equal(t, "c *:* rmw", dev.GetCgroupRule())
		default:
			assert.Fail(t, "unexpected device config")
		}
	}
}

func TestDeviceMappingsWithList(t *testing.T) {
	cfg, err := config.GetAllBackdrops("test/dodo.yaml")
	require.NoError(t, err)

	backdrop, ok := cfg["test_with_lists"]

	assert.True(t, ok)
	assert.Len(t, backdrop.GetDevices(), 1)

	for _, dev := range backdrop.GetDevices() {
		switch dev.GetTarget() {
		case "/dev/snd":
			assert.Equal(t, "/dev/snd", dev.GetSource())
		default:
			assert.Fail(t, "unexpected device config")
		}
	}
}
