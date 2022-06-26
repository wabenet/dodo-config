package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wabenet/dodo-config/internal/config"
)

func TestAllDefaults(t *testing.T) {
	cfg, err := config.GetAllBackdrops("test/dodo.yaml")

	assert.Nil(t, err)

	_, ok := cfg["test_all_defaults"]
	assert.True(t, ok)
}

func TestMinus(t *testing.T) {
	cfg, err := config.GetAllBackdrops("test/dodo.yaml")

	assert.Nil(t, err)

	_, ok := cfg["test-minus"]
	assert.True(t, ok)
}

func TestBasicBackdrop(t *testing.T) {
	cfg, err := config.GetAllBackdrops("test/dodo.yaml")

	assert.Nil(t, err)

	backdrop, ok := cfg["test_full_configs"]

	assert.True(t, ok)
	assert.Equal(t, "testimage", backdrop.ImageId)
	assert.Equal(t, "testcontainer", backdrop.ContainerName)
	assert.Equal(t, "/home/test", backdrop.WorkingDir)
	assert.Equal(t, "echo \"$@\"\n", backdrop.Entrypoint.Script)
}
