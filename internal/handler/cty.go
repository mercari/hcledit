package handler

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"

	"go.mercari.io/hcledit/internal/ast"
)

type ctyValueHandler struct {
	exprTokens    hclwrite.Tokens
	commentTokens hclwrite.Tokens

	afterKey string
}

func newCtyValueHandler(value cty.Value, comment, afterKey string) (Handler, error) {
	return &ctyValueHandler{
		exprTokens:    hclwrite.NewExpressionLiteral(value).BuildTokens(nil),
		commentTokens: commentTokens(comment),

		afterKey: afterKey,
	}, nil
}

func (h *ctyValueHandler) HandleObject(object *ast.Object, name string, _ []string) error {
	object.SetObjectAttributeRaw(name, h.exprTokens, h.commentTokens)
	if h.afterKey != "" {
		object.UpdateObjectAttributeOrder(name, h.afterKey)
	}
	return nil
}

func (h *ctyValueHandler) HandleBody(body *hclwrite.Body, name string, _ []string) error {
	body.SetAttributeRaw(name, h.exprTokens)

	if len(h.commentTokens) > 0 {
		tokens := body.GetAttribute(name).BuildTokens(h.commentTokens)
		body.RemoveAttribute(name)
		body.AppendUnstructuredTokens(tokens)
	}

	if h.afterKey != "" {
		ast.UpdateBodyTokenOrder(body, name, h.afterKey)
	}

	return nil
}
