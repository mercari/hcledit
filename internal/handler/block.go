package handler

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"

	"go.mercari.io/hcledit/internal/ast"
)

type BlockVal struct {
	Labels []string
}

type blockHandler struct {
	labels []string
}

func newBlockHandler(labels []string) (Handler, error) {
	return &blockHandler{
		labels: labels,
	}, nil
}

func (h *blockHandler) HandleBody(body *hclwrite.Body, name string, _ []string) error {
	body.AppendNewBlock(name, h.labels)
	return nil
}

func (h *blockHandler) HandleObject(_ *ast.Object, _ string, _ []string) error {
	return fmt.Errorf("this function should not be called")
}
