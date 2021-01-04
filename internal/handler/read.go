package handler

import (
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"

	"github.com/mercari/hcledit/internal/ast"
)

type readHandler struct {
	results map[string]cty.Value
}

func NewReadHandler(results map[string]cty.Value) (Handler, error) {
	return &readHandler{
		results: results,
	}, nil
}

func (h *readHandler) HandleBody(body *hclwrite.Body, name string, keyTrail []string) error {
	buf := body.GetAttribute(name).BuildTokens(nil).Bytes()
	value, err := parse(buf, name)
	if err != nil {
		return err
	}
	h.results[strings.Join(keyTrail, ".")] = value
	return nil
}

func (h *readHandler) HandleObject(object *ast.Object, name string, keyTrail []string) error {
	buf := object.GetObjectAttribute(name).BuildTokens().Bytes()
	value, err := parse(buf, name)
	if err != nil {
		return err
	}
	h.results[strings.Join(keyTrail, ".")] = value
	return nil
}

func parse(buf []byte, name string) (cty.Value, error) {
	file, diags := hclsyntax.ParseConfig(buf, "", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return cty.Value{}, diags
	}

	body := file.Body.(*hclsyntax.Body)
	v, diags := body.Attributes[name].Expr.Value(nil)
	if diags.HasErrors() {
		return cty.Value{}, diags
	}
	return v, nil
}
