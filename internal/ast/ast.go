package ast

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func ReplaceBodyTokens(body *hclwrite.Body, targetKey string, targetTokens hclwrite.Tokens) {
	originalKeys, newLines := getKeysAndNewLines(body.BuildTokens(nil))

	var (
		tokens  hclwrite.Tokens
		matched bool
	)
	for _, originalKey := range originalKeys {
		n := newLines[originalKey]
		for i := 0; i < n; i++ {
			tokens = append(tokens, hclwrite.Tokens{
				{
					Type:  hclsyntax.TokenNewline,
					Bytes: []byte{'\n'},
				},
			}...)
		}

		if originalKey == targetKey {
			tokens = append(tokens, targetTokens...)
			matched = true
			continue
		}
		tokens = append(tokens, body.GetAttribute(originalKey).BuildTokens(nil)...)
	}

	if matched {
		body.Clear()
		body.AppendUnstructuredTokens(tokens)
	}
}

func UpdateBodyTokenOrder(body *hclwrite.Body, targetKey string, afterKey string) {
	originalKeys, newLines := getKeysAndNewLines(body.BuildTokens(nil))

	var (
		tokens  hclwrite.Tokens
		matched bool
	)
	for _, originalKey := range originalKeys {
		if originalKey == targetKey {
			continue
		}
		n := newLines[originalKey]
		for i := 0; i < n; i++ {
			tokens = append(tokens, hclwrite.Tokens{
				{
					Type:  hclsyntax.TokenNewline,
					Bytes: []byte{'\n'},
				},
			}...)
		}

		tokens = append(tokens, getAttributeOrBlock(body, originalKey)...)
		if originalKey == afterKey {
			tokens = append(tokens, getAttributeOrBlock(body, targetKey)...)
			matched = true
		}
	}

	// Update order only when it finds after key in the original keys.
	if matched {
		body.Clear()
		body.AppendUnstructuredTokens(tokens)
	}
}

func getKeysAndNewLines(tokens hclwrite.Tokens) ([]string, map[string]int) {
	keys := []string{}
	newLineMap := make(map[string]int)

	p := &peeker{
		tokens: tokens,
	}

	var (
		open     []hclsyntax.TokenType
		newLines int
	)
Token:
	for {
		next := p.Peek()
		if next.Type == hclsyntax.TokenEOF {
			break Token
		}
		switch next.Type {
		case hclsyntax.TokenIdent:
			tokenIdent := p.Read()

			switch next := p.Peek(); next.Type {
			case hclsyntax.TokenEqual:
				// Attribute
				readAttribute(p)
			case hclsyntax.TokenOQuote, hclsyntax.TokenOBrace, hclsyntax.TokenIdent:
				// Block
				readBlock(p)
			default:
				continue
			}

			key := string(tokenIdent.Bytes)
			keys = append(keys, key)
			newLineMap[key] = newLines
			newLines = 0
		default:
			t := p.Read()
			if len(open) == 0 && t.Type == hclsyntax.TokenNewline {
				newLines++
			}
		}
	}
	return keys, newLineMap
}

func getAttributeOrBlock(body *hclwrite.Body, name string) hclwrite.Tokens {
	attr := body.GetAttribute(name)
	if attr != nil {
		return attr.BuildTokens(nil)
	}

	block := getBlock(body, name)
	return block.BuildTokens(nil)
}

func getBlock(body *hclwrite.Body, name string) *hclwrite.Block {
	for _, block := range body.Blocks() {
		if block.Type() == name {
			return block
		}
	}
	return nil
}
