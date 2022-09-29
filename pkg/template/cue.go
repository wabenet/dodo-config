package template

import (
	"reflect"

	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/literal"
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
	//
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
	// What is happening here might be a bit weird
	// So we check if the value is a quoted string, if yes we remove the
	// quotes and re-quote the string after templating. Reason is, the
	// re-quote should take care of escaping all newlines and so on that
	// come in through templating and would otherwise just destroy the AST.
	// We could also just always quote regardless, but that messes up
	// non-string literals (numbers, bools).
	// I realize this feels broken, and I'm sure there are use cases where
	// this will fall on our feet, but it needs someone smarter than me
	// right now to figure out how to handle this properly. I still have
	// hope the smarter person might be future-me eventually.
	//
	value := node.Value
	wasQuoted := false

	if unquoted, err := literal.Unquote(value); err == nil {
		value = unquoted
		wasQuoted = true
	}

	templated, err := ctx.TemplateString(value)
	if err != nil {
		return err
	}

	if wasQuoted {
		templated = literal.String.Quote(templated)
	}

	node.Value = templated

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
