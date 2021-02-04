package hcledit

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"

	"go.mercari.io/hcledit/internal/converter"
	"go.mercari.io/hcledit/internal/handler"
	"go.mercari.io/hcledit/internal/query"
	"go.mercari.io/hcledit/internal/walker"
)

// HCLEditor provides the interface of HCL editing.
type HCLEditor interface {
	// Create creates attributes and blocks matched with the given query
	// with the given value. The value can be any type and it's transformed
	// into HCL type inside.
	Create(query string, value interface{}, opts ...Option) error

	// Read returns attributes and blocks matched with the given query.
	// The results are map of mached key and its value.
	//
	// It returns error if it does not match any key.
	Read(query string, opts ...Option) (map[string]interface{}, error)

	// Update replaces attributes and blocks which matched with its key
	// and given query with the given value. The value can be any type
	// and it's transformed into HCL type inside.
	//
	// By default, it returns error if the does not matched with any key.
	// You must create value before update.
	Update(query string, value interface{}, opts ...Option) error

	// Delete deletes attributes and blocks matched with the given query.
	//
	// It returns error if it does not match any key.
	Delete(query string, opts ...Option) error

	Writer
}

type hclEditImpl struct {
	path      string
	filename  string
	writeFile *hclwrite.File
}

func BlockVal(labels ...string) *handler.BlockVal {
	return &handler.BlockVal{
		Labels: labels,
	}
}

func RawVal(rawString string) *handler.RawVal {
	return &handler.RawVal{
		RawString: rawString,
	}
}

func (h *hclEditImpl) Create(queryStr string, value interface{}, opts ...Option) error {
	defer h.reload()

	opt := &option{}
	for _, optFunc := range opts {
		optFunc(opt)
	}

	queries, err := query.Build(queryStr)
	if err != nil {
		return err
	}

	hdlr, err := handler.New(value, opt.comment, opt.afterKey)
	if err != nil {
		return err
	}

	w := &walker.Walker{
		Handler: hdlr,
		Mode:    walker.Create,
	}
	w.Walk(h.writeFile.Body(), queries, 0, []string{})

	return nil
}

func (h *hclEditImpl) Read(queryStr string, opts ...Option) (map[string]interface{}, error) {
	defer h.reload()

	opt := &option{}
	for _, optFunc := range opts {
		optFunc(opt)
	}

	queries, err := query.Build(queryStr)
	if err != nil {
		return nil, err
	}

	results := make(map[string]cty.Value)
	hdlr, err := handler.NewReadHandler(results)
	if err != nil {
		return nil, err
	}

	w := &walker.Walker{
		Handler: hdlr,
		Mode:    walker.Read,
	}

	w.Walk(h.writeFile.Body(), queries, 0, []string{})
	return convert(results)
}

func (h *hclEditImpl) Update(queryStr string, value interface{}, opts ...Option) error {
	defer h.reload()

	opt := &option{}
	for _, optFunc := range opts {
		optFunc(opt)
	}

	// TODO(tcnksm): Currently, WithComment is only for Create function.
	// We have the following challenges:
	// - We can not update existing comment. New comments are added
	//   top on the existing comments...
	// - Because we use `AppendUnstructuredTokens`and `RemoveAttribute`,
	//   the position of target attribute is also changed...
	if opt.comment != "" {
		return fmt.Errorf("WithComment is not supported for Update")
	}

	queries, err := query.Build(queryStr)
	if err != nil {
		return err
	}

	hdlr, err := handler.New(value, opt.comment, opt.afterKey)
	if err != nil {
		return err
	}

	w := &walker.Walker{
		Handler: hdlr,
		Mode:    walker.Update,
	}

	w.Walk(h.writeFile.Body(), queries, 0, []string{})
	return nil
}

func (h *hclEditImpl) Delete(queryStr string, opts ...Option) error {
	defer h.reload()

	opt := &option{}
	for _, optFunc := range opts {
		optFunc(opt)
	}

	queries, err := query.Build(queryStr)
	if err != nil {
		return err
	}

	w := &walker.Walker{
		Mode: walker.Delete,
	}

	w.Walk(h.writeFile.Body(), queries, 0, []string{})
	return nil
}

// reload re-parse the HCL file. Some operation causes like `WithAfter` modifies Body token structure
// drastically (it re-construct it from scratch...) and, because of it, some operation will not work
// properly after it
func (h *hclEditImpl) reload() error {
	writeFile, diags := hclwrite.ParseConfig(h.Bytes(), h.filename, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return diags
	}
	h.writeFile = writeFile
	return nil
}

func convert(original map[string]cty.Value) (map[string]interface{}, error) {
	results := make(map[string]interface{})
	for key, ctyVal := range original {
		goVal, err := converter.FromCtyValueToGoValue(ctyVal)
		if err != nil {
			return nil, err
		}
		results[key] = goVal
	}
	return results, nil
}
