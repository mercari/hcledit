// Package hcledit is a Go package to edit HCL configurations.
// Basically, this is just a wrapper of hclwrite package which provides
// low-level features of generating HCL configurations. But hcledit allows you
// to access HCL attribute or block by jq-like query and do various
// manipulations.
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

// HCLEditor implements an editor of HCL configuration.
type HCLEditor struct {
	path      string
	filename  string
	writeFile *hclwrite.File
}

// TODO(slewiskelly): Should these be exported?
// Users of this package are not allowed to import
// "go.mercari.io/hcledit/internal/handler" due to visibility of internal
// packages.
func BlockVal(labels ...string) *handler.BlockVal {
	return &handler.BlockVal{
		Labels: labels,
	}
}

// TODO(slewiskelly): As above.
func RawVal(rawString string) *handler.RawVal {
	return &handler.RawVal{
		RawString: rawString,
	}
}

// Create creates attributes and blocks matched with the given query
// with the given value. The value can be any type and it's transformed
// into HCL type inside.
func (h *HCLEditor) Create(queryStr string, value interface{}, opts ...Option) error {
	defer h.reload()

	opt := &option{
		querySeparator: '.',
	}
	for _, optFunc := range opts {
		optFunc(opt)
	}

	queries, err := query.Build(queryStr, opt.querySeparator)
	if err != nil {
		return err
	}

	hdlr, err := handler.New(value, opt.comment, opt.afterKey, opt.beforeNewline)
	if err != nil {
		return err
	}

	w := &walker.Walker{
		Handler: hdlr,
		Mode:    walker.Create,
	}

	return w.Walk(h.writeFile.Body(), queries, 0, []string{})
}

// Read returns attributes and blocks matched with the given query.
// The results are map of mached key and its value.
//
// It returns error if it does not match any key.
func (h *HCLEditor) Read(queryStr string, opts ...Option) (map[string]interface{}, error) {
	defer h.reload()

	opt := &option{
		querySeparator: '.',
	}
	for _, optFunc := range opts {
		optFunc(opt)
	}

	queries, err := query.Build(queryStr, opt.querySeparator)
	if err != nil {
		return nil, err
	}

	fallback := opt.readFallbackToRawString
	results := make(map[string]cty.Value)
	hdlr, err := handler.NewReadHandler(results, fallback)
	if err != nil {
		return nil, err
	}

	w := &walker.Walker{
		Handler: hdlr,
		Mode:    walker.Read,
	}

	walkErr := w.Walk(h.writeFile.Body(), queries, 0, []string{})
	if walkErr != nil && !fallback {
		return nil, walkErr
	}

	ret, convertErr := convert(results)
	if convertErr != nil {
		return ret, convertErr
	}

	return ret, walkErr
}

// Update replaces attributes and blocks which matched with its key
// and given query with the given value. The value can be any type
// and it's transformed into HCL type inside.
//
// By default, it returns error if the does not matched with any key.
// You must create value before update.
func (h *HCLEditor) Update(queryStr string, value interface{}, opts ...Option) error {
	defer h.reload()

	opt := &option{
		querySeparator: '.',
	}
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
	} else if opt.beforeNewline {
		return fmt.Errorf("WithNewLine is not supported for Update")
	}

	queries, err := query.Build(queryStr, opt.querySeparator)
	if err != nil {
		return err
	}

	hdlr, err := handler.New(value, opt.comment, opt.afterKey, opt.beforeNewline)
	if err != nil {
		return err
	}

	w := &walker.Walker{
		Handler: hdlr,
		Mode:    walker.Update,
	}

	return w.Walk(h.writeFile.Body(), queries, 0, []string{})
}

// Delete deletes attributes and blocks matched with the given query.
//
// It returns error if it does not match any key.
func (h *HCLEditor) Delete(queryStr string, opts ...Option) error {
	defer h.reload()

	opt := &option{
		querySeparator: '.',
	}
	for _, optFunc := range opts {
		optFunc(opt)
	}

	queries, err := query.Build(queryStr, opt.querySeparator)
	if err != nil {
		return err
	}

	w := &walker.Walker{
		Mode: walker.Delete,
	}

	return w.Walk(h.writeFile.Body(), queries, 0, []string{})
}

// CustomEdit executes a custom function on the underlying file
func (h *HCLEditor) CustomEdit(fn func(*hclwrite.Body) error) error {
	defer h.reload()
	return fn(h.writeFile.Body())
}

// reload re-parse the HCL file. Some operation causes like `WithAfter` modifies Body token structure
// drastically (it re-construct it from scratch...) and, because of it, some operation will not work
// properly after it
func (h *HCLEditor) reload() error {
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
