package handler

import (
	"fmt"
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
	attr := body.GetAttribute(name)
	if attr == nil {
		return fmt.Errorf("attribute %s not found", name)
	}

	buf := attr.BuildTokens(nil).Bytes()
	fallback := h.fallbackToRawString
	value, err := parse(buf, name, fallback)
	if err != nil && !fallback {
		return err
	}
	h.results[strings.Join(keyTrail, ".")] = value
	return err
}

func (h *readHandler) HandleObject(object *ast.Object, name string, keyTrail []string) error {
	attr := object.GetObjectAttribute(name)
	if attr == nil {
		return fmt.Errorf("attribute %s not found", name)
	}

	buf := attr.BuildTokens().Bytes()
	fallback := h.fallbackToRawString
	value, err := parse(buf, name, fallback)
	if err != nil && !fallback {
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
	attr, ok := body.Attributes[name]
	if !ok {
		return cty.Value{}, fmt.Errorf("attribute %s not found", name)
	}

	v, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		if !fallback {
			return cty.Value{}, diags
		}

		// Fallback: Return raw string representation of the object
		rawValue := string(attr.Expr.Range().SliceBytes(buf))
		return cty.StringVal(rawValue), nil
	}

	return v, nil
}
