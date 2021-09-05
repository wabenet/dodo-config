package config_test

import (
	"fmt"
	"testing"

	"github.com/dodo-cli/dodo-config/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestAllDefaults(t *testing.T) {
	config, err := config.ParseConfig("test/dodo.yaml")
	if err != nil {
		fmt.Printf(err.Error())
	} else {
		fmt.Printf("%v\n", config)
	}

	assert.Nil(t, err)

	_, ok := config["test_all_defaults"]
	assert.True(t, ok)
}

func TestBasicBackdrop(t *testing.T) {
	config, err := config.ParseConfig("test/dodo.yaml")

	assert.Nil(t, err)

	backdrop, ok := config["test_full_configs"]

	assert.True(t, ok)
	assert.Equal(t, "testimage", backdrop.ImageId)
	assert.Equal(t, "testcontainer", backdrop.ContainerName)
	assert.Contains(t, backdrop.Entrypoint.Interpreter, "/bin/sh")
	assert.Equal(t, "/home/test", backdrop.WorkingDir)
	assert.Equal(t, "echo \"$@\"\n", backdrop.Entrypoint.Script)
}
