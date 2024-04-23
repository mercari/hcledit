package handler

import (
	"strings"
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"

	"go.mercari.io/hcledit/internal/ast"
)

type BlockVal struct {
	Labels []string
}

type blockHandler struct {
	labels  []string
	comment string
}

func newBlockHandler(labels []string, comment string) (Handler, error) {
	return &blockHandler{
		labels:  labels,
		comment: comment,
	}, nil
}

func (h *blockHandler) HandleBody(body *hclwrite.Body, name string, _ []string) error {
	if h.comment != "" {
		if !strings.HasPrefix(h.comment, "//") {
			h.comment = fmt.Sprintf("// %s", h.comment)
		}
		tokens := beforeTokens(h.comment, false)
		body.AppendUnstructuredTokens(tokens)
	}

	body.AppendNewBlock(name, h.labels)
	return nil
}

func (h *blockHandler) HandleObject(_ *ast.Object, _ string, _ []string) error {
	return fmt.Errorf("this function should not be called")
}
