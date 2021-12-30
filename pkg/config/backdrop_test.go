package config_test

import (
	"testing"

	"github.com/dodo-cli/dodo-config/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestAllDefaults(t *testing.T) {
	cfg, err := config.ParseConfig("test/dodo.yaml")

	assert.Nil(t, err)

	_, ok := cfg.Backdrops["test_all_defaults"]
	assert.True(t, ok)
}

func TestMinus(t *testing.T) {
	cfg, err := config.ParseConfig("test/dodo.yaml")

	assert.Nil(t, err)

	_, ok := cfg.Backdrops["test-minus"]
	assert.True(t, ok)
}


func TestBasicBackdrop(t *testing.T) {
	cfg, err := config.ParseConfig("test/dodo.yaml")

	assert.Nil(t, err)

	backdrop, ok := cfg.Backdrops["test_full_configs"]

	assert.True(t, ok)
	assert.Equal(t, "testimage", backdrop.ImageId)
	assert.Equal(t, "testcontainer", backdrop.ContainerName)
	assert.Equal(t, "/home/test", backdrop.WorkingDir)
	assert.Equal(t, "echo \"$@\"\n", backdrop.Entrypoint.Script)
}
