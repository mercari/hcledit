package handler

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"

	"go.mercari.io/hcledit/internal/ast"
)

type ctyValueHandler struct {
	exprTokens   hclwrite.Tokens
	beforeTokens hclwrite.Tokens

	afterKey string
}

func newCtyValueHandler(value cty.Value, comment, afterKey string, beforeNewline bool) (Handler, error) {
	return &ctyValueHandler{
		exprTokens:   hclwrite.NewExpressionLiteral(value).BuildTokens(nil),
		beforeTokens: beforeTokens(comment, beforeNewline),

		afterKey: afterKey,
	}, nil
}

func (h *ctyValueHandler) HandleObject(object *ast.Object, name string, _ []string) error {
	object.SetObjectAttributeRaw(name, h.exprTokens, h.beforeTokens)
	if h.afterKey != "" {
		object.UpdateObjectAttributeOrder(name, h.afterKey)
	}
	return nil
}

func (h *ctyValueHandler) HandleBody(body *hclwrite.Body, name string, _ []string) error {
	body.SetAttributeRaw(name, h.exprTokens)

	if len(h.beforeTokens) > 0 {
		tokens := body.GetAttribute(name).BuildTokens(h.beforeTokens)
		body.RemoveAttribute(name)
		body.AppendUnstructuredTokens(tokens)
	}

	if h.afterKey != "" {
		ast.UpdateBodyTokenOrder(body, name, h.afterKey)
	}

	return nil
}
