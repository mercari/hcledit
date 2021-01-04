package handler

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty/gocty"

	"github.com/mercari/hcledit/internal/ast"
)

type Handler interface {
	HandleBody(body *hclwrite.Body, name string, keyTrail []string) error
	HandleObject(object *ast.Object, name string, keyTrail []string) error
}

func New(input interface{}, comment, afterKey string) (Handler, error) {
	switch v := input.(type) {
	case *BlockVal:
		return newBlockHandler(v.Labels)
	}

	ctyType, err := gocty.ImpliedType(input)
	if err != nil {
		return nil, fmt.Errorf("failed to imply cty type from go value: %s", err)
	}

	ctyVal, err := gocty.ToCtyValue(input, ctyType)
	if err != nil {
		return nil, fmt.Errorf("failed to convert cty value from go value: %s", err)
	}

	return newCtyValueHandler(ctyVal, comment, afterKey)
}

func commentTokens(comment string) hclwrite.Tokens {
	if comment == "" {
		return hclwrite.Tokens{}
	}

	return hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenComment,
			Bytes: []byte(comment),
		},
		{
			Type:  hclsyntax.TokenNewline,
			Bytes: []byte{'\n'},
		},
	}
}
