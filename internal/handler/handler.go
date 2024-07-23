package handler

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty/gocty"

	"go.mercari.io/hcledit/internal/ast"
)

type Handler interface {
	HandleBody(body *hclwrite.Body, name string, keyTrail []string) error
	HandleObject(object *ast.Object, name string, keyTrail []string) error
}

func New(input interface{}, comment, afterKey string, beforeNewline bool) (Handler, error) {
	switch v := input.(type) {
	case *BlockVal:
		return newBlockHandler(v.Labels, comment, beforeNewline)
	case *RawVal:
		return newRawHandler(v.RawString, comment, beforeNewline)
	}

	ctyType, err := gocty.ImpliedType(input)
	if err != nil {
		return nil, fmt.Errorf("failed to imply cty type from go value: %s", err)
	}

	ctyVal, err := gocty.ToCtyValue(input, ctyType)
	if err != nil {
		return nil, fmt.Errorf("failed to convert cty value from go value: %s", err)
	}

	return newCtyValueHandler(ctyVal, comment, afterKey, beforeNewline)
}

func beforeTokens(comment string, beforeNewline bool) hclwrite.Tokens {
	result := hclwrite.Tokens{}
	if beforeNewline {
		result = append(result, &hclwrite.Token{
			Type:  hclsyntax.TokenNewline,
			Bytes: []byte{'\n'},
		})
	}

	if comment != "" {
		result = append(result, &hclwrite.Token{
			Type:  hclsyntax.TokenComment,
			Bytes: []byte(comment),
		})
		result = append(result, &hclwrite.Token{
			Type:  hclsyntax.TokenNewline,
			Bytes: []byte{'\n'},
		})
	}

	return result
}
