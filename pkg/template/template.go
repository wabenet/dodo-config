package template

import (
	"bytes"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig"
)

type Context struct {
	Filename string
}

func (ctx *Context) TemplateString(input string) (string, error) {
	templ, err := template.New("config").Funcs(sprig.TxtFuncMap()).Funcs(ctx.FuncMap()).Parse(input)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer

	if err := templ.Execute(&buffer, nil); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func (ctx *Context) FuncMap() template.FuncMap {
	return template.FuncMap{
		"user":        user.Current,
		"cwd":         os.Getwd,
		"env":         os.Getenv,
		"currentFile": ctx.currentFile,
		"currentDir":  ctx.currentDir,
		"sh":          ctx.runShell,
		"readFile":    ctx.readFile,
		"projectRoot": ctx.projectRoot,
		"projectPath": ctx.projectPath,
	}
}

func (ctx *Context) currentFile() string {
	return ctx.Filename
}

func (ctx *Context) currentDir() string {
	return filepath.Dir(ctx.Filename)
}

func (*Context) runShell(command string) (string, error) {
	// TODO: what to do on windows?
	var out bytes.Buffer

	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return out.String(), nil
}

func (ctx *Context) readFile(path string) (string, error) {
	fullPath := filepath.Join(ctx.currentDir(), path)

	b, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (*Context) projectRoot() (string, error) {
	root, _, err := findProjectRoot()

	return root, err
}

func (*Context) projectPath() (string, error) {
	_, path, err := findProjectRoot()

	return path, err
}

func findProjectRoot() (string, string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", "", err
	}

	for dir := cwd; dir != "/"; dir = filepath.Dir(dir) {
		if info, err := os.Stat(filepath.Join(dir, ".git")); err == nil && info.IsDir() {
			path, err2 := filepath.Rel(dir, cwd)
			if err2 != nil {
				return "", "", err
			}

			return dir, path, nil
		}
	}

	return cwd, ".", nil
}
