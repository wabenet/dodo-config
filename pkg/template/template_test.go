package template_test

import (
	"os"
	"path/filepath"
	"testing"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/wabenet/dodo-config/pkg/template"
)

func TestTemplateYamlFile(t *testing.T) {
	t.Parallel()

	yamlFile, err := yaml.Extract("./test/test.yaml", nil)
	assert.Nil(t, err)

	yamlFile, err = template.TemplateCueAST(yamlFile)
	assert.Nil(t, err)

	ctx := cuecontext.New()
	bi := load.Instances([]string{"-"}, nil)[0]

	err = bi.AddSyntax(yamlFile)
	assert.Nil(t, err)

	value := ctx.BuildInstance(bi)
	assert.Nil(t, value.Err())

	err = value.Validate(cue.Concrete(true), cue.Final())
	assert.Nil(t, err)

	cwd := get(t, value, "cwd")
	currentDir := get(t, value, "currentDir")
	currentFile := get(t, value, "currentFile")
	projectRoot := get(t, value, "projectRoot")
	projectPath := get(t, value, "projectPath")

	assert.Equal(t, "pkg/template", projectPath)
	assert.Equal(t, filepath.Join(projectRoot, "pkg/template"), cwd)
	// TODO: these should always be absolute!
	assert.Equal(t, "test", currentDir)
	assert.Equal(t, "./test/test.yaml", currentFile)

	hello := get(t, value, "hello")
	assert.Equal(t, "Hello World\n", hello)

	included, err := os.ReadFile("./test/included.txt")
	assert.Nil(t, err)

	readfile := get(t, value, "readfile")
	assert.Equal(t, string(included), readfile)

	readfileMultiline := get(t, value, "readfile_multiline")
	assert.Equal(t, string(included), readfileMultiline)
}

func get(t *testing.T, v cue.Value, name string) string {
	t.Helper()

	s, err := v.LookupPath(cue.MakePath(cue.Str(name))).String()
	assert.Nil(t, err)

	return s
}
