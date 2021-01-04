package ast

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type peeker struct {
	tokens    hclwrite.Tokens
	nextIndex int
}

func (r *peeker) Peek() *hclwrite.Token {
	ret, _ := r.nextToken()
	return ret
}

func (r *peeker) Read() *hclwrite.Token {
	token, nextIndex := r.nextToken()
	r.nextIndex = nextIndex
	return token
}

func (r *peeker) nextToken() (*hclwrite.Token, int) {
	if r.nextIndex > len(r.tokens)-1 {
		return &hclwrite.Token{
			Type: hclsyntax.TokenEOF,
		}, len(r.tokens)
	}
	return r.tokens[r.nextIndex], r.nextIndex + 1
}
