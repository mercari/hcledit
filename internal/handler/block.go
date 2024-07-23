package handler

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"

	"go.mercari.io/hcledit/internal/ast"
)

type BlockVal struct {
	Labels []string
}

type blockHandler struct {
	labels        []string
	comment       string
	beforeNewLine bool
}

func newBlockHandler(labels []string, comment string, beforeNewLine bool) (Handler, error) {
	return &blockHandler{
		labels:        labels,
		comment:       comment,
		beforeNewLine: beforeNewLine,
	}, nil
}

func (h *blockHandler) HandleBody(body *hclwrite.Body, name string, _ []string) error {
	if h.comment != "" {
		// Note: intuitively, adding a new line should be determined by `h.beforeNewLine`.
		// However, for the backward compatibility, we add a new line whenever comment is set to a non empty string.
		body.AppendUnstructuredTokens(
			beforeTokens(
				fmt.Sprintf("// %s", strings.TrimSpace(strings.TrimPrefix(h.comment, "//"))),
				true,
			),
		)
	} else {
		body.AppendUnstructuredTokens(beforeTokens(h.comment, h.beforeNewLine))
	}
	body.AppendNewBlock(name, h.labels)
	return nil
}

func (h *blockHandler) HandleObject(_ *ast.Object, _ string, _ []string) error {
	return fmt.Errorf("this function should not be called")
}
