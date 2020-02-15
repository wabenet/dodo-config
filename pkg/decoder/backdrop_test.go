package decoder

import (
	"testing"

	"github.com/oclaussen/dodo/pkg/types"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

const simplePull = `
image: testimage
`

func TestSimplePull(t *testing.T) {
	config := getExampleConfig(t, simplePull)
	assert.Equal(t, "testimage", config.ImageId)
}

const environments = `
environment:
  - FOO=BAR
  - SOMETHING
`

func TestEnvironments(t *testing.T) {
	config := getExampleConfig(t, environments)
	assert.Equal(t, 2, len(config.Environment))

	foobar := config.Environment[0]
	assert.Equal(t, foobar.Key, "FOO")
	assert.NotNil(t, foobar.Value)
	assert.Equal(t, "BAR", foobar.Value)

	something := config.Environment[1]
	assert.Equal(t, "SOMETHING", something.Key)
	assert.Equal(t, "", something.Value)
}

const simpleVolume = `
volumes: foo:bar:ro
`

func TestSimpleVolume(t *testing.T) {
	config := getExampleConfig(t, simpleVolume)
	assert.Equal(t, 1, len(config.Volumes))
	assert.Equal(t, "foo", config.Volumes[0].Source)
	assert.Equal(t, "bar", config.Volumes[0].Target)
	assert.True(t, config.Volumes[0].Readonly)
}

const mixedVolumes = `
volumes:
  - test
  - source: from
    target: to
    read_only: true
  - bar:baz
`

func TestMixedVolumes(t *testing.T) {
	config := getExampleConfig(t, mixedVolumes)
	assert.Equal(t, 3, len(config.Volumes))

	sourceOnly := config.Volumes[0]
	assert.Equal(t, "test", sourceOnly.Source)
	assert.Equal(t, "", sourceOnly.Target)
	assert.False(t, sourceOnly.Readonly)

	fullSpec := config.Volumes[1]
	assert.Equal(t, "from", fullSpec.Source)
	assert.Equal(t, "to", fullSpec.Target)
	assert.True(t, fullSpec.Readonly)

	readWrite := config.Volumes[2]
	assert.Equal(t, "bar", readWrite.Source)
	assert.Equal(t, "baz", readWrite.Target)
	assert.False(t, readWrite.Readonly)
}

const fullExample = `
image: testimage
container_name: testcontainer
remove: false
interactive: true
volumes_from: 'somevolume'
interpreter: '/bin/sh'
script: |
  echo "$@"
command: ['Hello', 'World']
`

func TestFullExample(t *testing.T) {
	config := getExampleConfig(t, fullExample)
	assert.Equal(t, "testimage", config.ImageId)
	assert.Equal(t, "testcontainer", config.ContainerName)
	assert.True(t, config.Entrypoint.Interactive)
	assert.Contains(t, config.Entrypoint.Interpreter, "/bin/sh")
	assert.Equal(t, "echo \"$@\"\n", config.Entrypoint.Script)
	assert.Equal(t, []string{"Hello", "World"}, config.Entrypoint.Arguments)
}

func getExampleConfig(t *testing.T, yamlConfig string) types.Backdrop {
	var mapType map[interface{}]interface{}
	err := yaml.Unmarshal([]byte(yamlConfig), &mapType)
	assert.Nil(t, err)
	decoder := NewDecoder("example")
	config, err := decoder.DecodeBackdrop("example", mapType)
	assert.Nil(t, err)
	return config
}