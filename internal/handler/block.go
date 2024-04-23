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
		body.AppendUnstructuredTokens(
			beforeTokens(
				fmt.Sprintf("// %s", strings.TrimSpace(strings.TrimPrefix(h.comment, "//"))),
				false,
			),
		)
	}

	body.AppendNewBlock(name, h.labels)
	return nil
}

func (h *blockHandler) HandleObject(_ *ast.Object, _ string, _ []string) error {
	return fmt.Errorf("this function should not be called")
}
