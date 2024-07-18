package handler

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"

	"go.mercari.io/hcledit/internal/ast"
)

type RawVal struct {
	RawString string
}

type rawHandler struct {
	rawTokens hclwrite.Tokens
	beforeTokens hclwrite.Tokens
}

func newRawHandler(rawString, comment string, beforeNewline bool) (Handler, error) {
	return &rawHandler{
		// NOTE(tcnksm):
		rawTokens: hclwrite.Tokens{
			{
				Type:  hclsyntax.TokenComment,
				Bytes: []byte(rawString),
			},
		},
		beforeTokens: beforeTokens(comment, beforeNewline),
	}, nil
}

func (h *rawHandler) HandleBody(body *hclwrite.Body, name string, _ []string) error {
	body.SetAttributeRaw(name, h.rawTokens)
	if len(h.beforeTokens) > 0 {
		tokens := body.GetAttribute(name).BuildTokens(h.beforeTokens)
		body.RemoveAttribute(name)
		body.AppendUnstructuredTokens(tokens)
	}
	return nil
}

func (h *rawHandler) HandleObject(object *ast.Object, name string, _ []string) error {
	object.SetObjectAttributeRaw(name, h.rawTokens, hclwrite.Tokens{})
	return nil
}
