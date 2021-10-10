package template

import (
	"reflect"

	"cuelang.org/go/cue/ast"
	log "github.com/hashicorp/go-hclog"
)

func TemplateCueAST(file *ast.File) (*ast.File, error) {
	ctx := &Context{Filename: file.Filename}

	for _, decl := range file.Decls {
		if err := ctx.templateCueNode(decl); err != nil {
			return nil, err
		}
	}

	return file, nil
}

func (ctx *Context) templateCueNode(node ast.Node) error {
	// In theory we could add every CUE AST type here
	// But that feels like it contradicts the idea of using CUE in the first place
	// So we only template those types that specifically occur in YAML

	if n, ok := node.(*ast.BasicLit); ok {
		return ctx.templateCueBasicLit(n)
	}

	if n, ok := node.(*ast.Field); ok {
		return ctx.templateCueField(n)
	}

	if n, ok := node.(*ast.ListLit); ok {
		return ctx.templateCueListLit(n)
	}

	if n, ok := node.(*ast.StructLit); ok {
		return ctx.templateCueStructLit(n)
	}

	log.L().Warn("Declaration ignored for templating", "decl", node, "type", reflect.TypeOf(node))
        return nil
}

func (ctx *Context) templateCueBasicLit(node *ast.BasicLit) error {
	v, err := ctx.TemplateString(node.Value)
	if err != nil {
		return err
	}

	node.Value = v
	return nil
}

func (ctx *Context) templateCueField(node *ast.Field) error {
	// Don't allow templating of the label, only the value
	return ctx.templateCueNode(node.Value)
}

func (ctx *Context) templateCueListLit(node *ast.ListLit) error {
	for _, elt := range node.Elts {
		if err := ctx.templateCueNode(elt); err != nil {
			return err
		}
	}

	return nil
}

func (ctx *Context) templateCueStructLit(node *ast.StructLit) error {
	for _, elt := range node.Elts {
		if err := ctx.templateCueNode(elt); err != nil {
			return err
		}
	}

	return nil
}
