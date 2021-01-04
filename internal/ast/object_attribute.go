package ast

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// ObjectAttribute is attribute inside Object
type ObjectAtrribute struct {
	name string

	beforeTokens hclwrite.Tokens
	exprTokens   hclwrite.Tokens
	afterTokens  hclwrite.Tokens
}

func (oa *ObjectAtrribute) BuildTokens() hclwrite.Tokens {
	var tokens hclwrite.Tokens
	tokens = append(tokens, oa.beforeTokens...)
	tokens = append(tokens, oa.exprTokens...)
	tokens = append(tokens, oa.afterTokens...)
	return tokens
}

func newObjectAttribute(name string, exprTokens, commentTokens hclwrite.Tokens) *ObjectAtrribute {
	beforeTokens := hclwrite.Tokens{}
	if len(commentTokens) != 0 {
		beforeTokens = append(beforeTokens, commentTokens...)
	}
	beforeTokens = append(beforeTokens, hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte(name),
		},
		{
			Type:  hclsyntax.TokenEqual,
			Bytes: []byte{'='},
		},
	}...)

	return &ObjectAtrribute{
		name: name,

		beforeTokens: beforeTokens,
		exprTokens:   exprTokens,
		afterTokens: hclwrite.Tokens{
			{
				Type:  hclsyntax.TokenNewline,
				Bytes: []byte{'\n'},
			},
		},
	}
}
