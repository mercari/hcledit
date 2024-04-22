package handler

import (
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"

	"go.mercari.io/hcledit/internal/ast"
)

type readHandler struct {
	results             map[string]cty.Value
	fallbackToRawString bool
}

func NewReadHandler(results map[string]cty.Value, fallbackToRawString bool) (Handler, error) {
	return &readHandler{
		results:             results,
		fallbackToRawString: fallbackToRawString,
	}, nil
}

func (h *readHandler) HandleBody(body *hclwrite.Body, name string, keyTrail []string) error {
	buf := body.GetAttribute(name).BuildTokens(nil).Bytes()
	value, err := parse(buf, name, h.fallbackToRawString)
	if err != nil {
		return err
	}
	h.results[strings.Join(keyTrail, ".")] = value
	return nil
}

func (h *readHandler) HandleObject(object *ast.Object, name string, keyTrail []string) error {
	buf := object.GetObjectAttribute(name).BuildTokens().Bytes()
	value, err := parse(buf, name, h.fallbackToRawString)
	if err != nil {
		return err
	}
	h.results[strings.Join(keyTrail, ".")] = value
	return nil
}

func parse(buf []byte, name string, fallback bool) (cty.Value, error) {
	file, diags := hclsyntax.ParseConfig(buf, "", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return cty.Value{}, diags
	}

	body := file.Body.(*hclsyntax.Body)
	expr := body.Attributes[name].Expr
	v, diags := expr.Value(nil)
	if diags.HasErrors() {
		if !fallback {
			return cty.Value{}, diags
		}

		// Could not parse the value with a nil EvalContext, so this is likely an
		// interpolated string. Instead, attempt to parse the raw string value.
		return cty.StringVal(string(expr.Range().SliceBytes(buf))), nil
	}
	return v, nil
}
