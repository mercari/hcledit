package ast

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func ParseObject(tokens hclwrite.Tokens) (*Object, error) {
	beforeObjTokens, exprObjTokens, afterObjTokens, err := partitionObjectTokens(tokens)
	if err != nil {
		return nil, err
	}

	obj := &Object{
		beforeTokens: beforeObjTokens,
		afterTokens:  afterObjTokens,
	}

	p := &peeker{
		tokens: exprObjTokens,
	}

	var beforeTokens hclwrite.Tokens
Token:
	for {
		next := p.Peek()
		if next.Type == hclsyntax.TokenEOF {
			break Token
		}

		switch next.Type {
		case hclsyntax.TokenIdent:
			tokenIdent := p.Read()
			beforeTokens = append(beforeTokens, tokenIdent)

			if next := p.Peek(); next.Type != hclsyntax.TokenEqual {
				continue
			}

			parsedBeforeTokens, exprTokens, afterTokens := readAttribute(p)
			beforeTokens = append(beforeTokens, parsedBeforeTokens...)

			obj.objectAtrributes = append(obj.objectAtrributes, &ObjectAtrribute{
				name:         string(tokenIdent.Bytes),
				beforeTokens: beforeTokens,
				exprTokens:   exprTokens,
				afterTokens:  afterTokens,
			})

			// Initialize for next read
			beforeTokens = hclwrite.Tokens{}
		default:
			// Handle all tokens 'before' TokenIdent as beforeTokens.
			// This should be mainly comments or new lines.
			beforeTokens = append(beforeTokens, p.Read())
		}
	}
	return obj, nil
}

func partitionObjectTokens(tokens hclwrite.Tokens) (hclwrite.Tokens, hclwrite.Tokens, hclwrite.Tokens, error) {
	var start int
	for i := 0; i <= len(tokens)-1; i++ {
		if tokens[i].Type == hclsyntax.TokenOBrace {
			// NOTE(tcnksm): Increment 1 to include new line
			start = i + 1
			break
		}
	}

	var end int
	for i := len(tokens) - 1; i >= 0; i-- {
		if tokens[i].Type == hclsyntax.TokenCBrace {
			end = i
			break
		}
	}
	if start == end {
		return nil, nil, nil, fmt.Errorf("invalid")
	}

	return tokens[0 : start+1], tokens[start+1 : end], tokens[end:], nil
}

// readBlock is only used to move peeker to end of the block expression.
func readBlock(p *peeker) {
BeforeToken:
	for {
		tok := p.Read()
		if tok.Type == hclsyntax.TokenEOF || tok.Type == hclsyntax.TokenOBrace {
			break BeforeToken
		}
	}

	var open []hclsyntax.TokenType
Token:
	for {
		next := p.Peek()
		if next.Type == hclsyntax.TokenEOF {
			break Token
		}
		switch next.Type {
		case hclsyntax.TokenOBrace:
			token := p.Read()
			open = append(open, token.Type)
		case hclsyntax.TokenCBrace:
			token := p.Read()
			if len(open) == 0 {
				p.Read() // eat newline Token
				break Token
			}

			opener := oppositeBracket(token.Type)
			for len(open) > 0 && open[len(open)-1] != opener {
				open = open[:len(open)-1]
			}
			if len(open) > 0 {
				open = open[:len(open)-1]
			}
		default:
			p.Read()
		}
	}
}

func readAttribute(p *peeker) (hclwrite.Tokens, hclwrite.Tokens, hclwrite.Tokens) {
	var (
		beforeTokens hclwrite.Tokens
		exprTokens   hclwrite.Tokens
		afterTokens  hclwrite.Tokens

		open []hclsyntax.TokenType
	)

	tokenEqual := p.Read()
	if tokenEqual.Type != hclsyntax.TokenEqual {
		panic("readAttribute must be called when it's attribute type")
	}
	beforeTokens = append(beforeTokens, tokenEqual)

EndExpr:
	for {
		next := p.Peek()
		if next.Type == hclsyntax.TokenEOF {
			break
		}

		switch next.Type {
		case hclsyntax.TokenOBrace, hclsyntax.TokenOBrack, hclsyntax.TokenOParen, hclsyntax.TokenOQuote, hclsyntax.TokenOHeredoc:
			token := p.Read()
			exprTokens = append(exprTokens, token)
			open = append(open, token.Type)
		case hclsyntax.TokenCBrace, hclsyntax.TokenCBrack, hclsyntax.TokenCParen, hclsyntax.TokenCQuote, hclsyntax.TokenCHeredoc:
			token := p.Read()
			opener := oppositeBracket(token.Type)
			for len(open) > 0 && open[len(open)-1] != opener {
				open = open[:len(open)-1]
			}
			if len(open) > 0 {
				open = open[:len(open)-1]
			}
			exprTokens = append(exprTokens, token)
		case hclsyntax.TokenComment:
			token := p.Read()

			// This should be line comment.
			if len(open) == 0 {
				afterTokens = append(afterTokens, token)
				break EndExpr
			}

			exprTokens = append(exprTokens, token)
		case hclsyntax.TokenNewline:
			token := p.Read()

			// This should be end of expression.
			if len(open) == 0 {
				afterTokens = append(afterTokens, token)
				break EndExpr
			}
			exprTokens = append(exprTokens, token)
		default:
			exprTokens = append(exprTokens, p.Read())
		}
	}

	return beforeTokens, exprTokens, afterTokens
}

// This is copied from github.com/hashicorp/hcl/v2@v2.8.1/hclsyntax/parser.go
func oppositeBracket(ty hclsyntax.TokenType) hclsyntax.TokenType {
	switch ty {
	case hclsyntax.TokenOBrace:
		return hclsyntax.TokenCBrace
	case hclsyntax.TokenOBrack:
		return hclsyntax.TokenCBrack
	case hclsyntax.TokenOParen:
		return hclsyntax.TokenCParen
	case hclsyntax.TokenOQuote:
		return hclsyntax.TokenCQuote
	case hclsyntax.TokenOHeredoc:
		return hclsyntax.TokenCHeredoc
	case hclsyntax.TokenCBrace:
		return hclsyntax.TokenOBrace
	case hclsyntax.TokenCBrack:
		return hclsyntax.TokenOBrack
	case hclsyntax.TokenCParen:
		return hclsyntax.TokenOParen
	case hclsyntax.TokenCQuote:
		return hclsyntax.TokenOQuote
	case hclsyntax.TokenCHeredoc:
		return hclsyntax.TokenOHeredoc
	default:
		return hclsyntax.TokenNil
	}
}
