package handler

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/mercari/hcledit/internal/ast"
)

type RawVal struct {
	RawString string
}

type rawHandler struct {
	rawTokens hclwrite.Tokens
}

func newRawHandler(rawString string) (Handler, error) {
	return &rawHandler{
		// NOTE(tcnksm):
		rawTokens: hclwrite.Tokens{
			{
				Type:  hclsyntax.TokenComment,
				Bytes: []byte(rawString),
			},
		}}, nil
}

func (h *rawHandler) HandleBody(body *hclwrite.Body, name string, _ []string) error {
	body.SetAttributeRaw(name, h.rawTokens)
	return nil
}

func (h *rawHandler) HandleObject(object *ast.Object, name string, _ []string) error {
	object.SetObjectAttributeRaw(name, h.rawTokens, hclwrite.Tokens{})
	return nil
}
