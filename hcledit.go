// Package hcledit is a Go package to edit HCL configurations.
// Basically, this is just a wrapper of hclwrite package which provides
// low-level features of generating HCL configurations. But hcledit allows you
// to access HCL attribute or block by jq-like query and do various
// manipulations.
package hcledit

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
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
	hclFile   *hcl.File
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

type root struct {
	Variables []*variable `hcl:"variable,block"`
	Locals    []*local    `hcl:"locals,block"`
}

type variable struct {
	Name    string         `hcl:"name,label"`
	Default hcl.Expression `hcl:"default,optional"`
	Remain  hcl.Body       `hcl:",remain"`
}

type local struct {
	Remain hcl.Body `hcl:",remain"`
}

func (h *HCLEditor) updateVariables(ctx *hcl.EvalContext, vars ...*variable) {
	const key = "var"
	children := make(map[string]cty.Value)
	for _, v := range vars {
		value, _ := v.Default.Value(ctx)

		children[v.Name] = value

		// TODO(ryan-ph): This needs to handle the case where key already
		// exists in this map. If it does exist, we need to merge the vals
		//
		// update context as other blocks or attrs can refer to values we just
		// parsed
		ctx.Variables[key] = cty.ObjectVal(children)
	}
}

func (h *HCLEditor) updateLocals(ctx *hcl.EvalContext, locals ...*local) {
	const key = "local"
	children := make(map[string]cty.Value)
	for _, l := range locals {
		attrs, _ := l.Remain.JustAttributes()
		for name, attr := range attrs {
			value, _ := attr.Expr.Value(ctx)

			children[name] = value

			// TODO(ryan-ph): This needs to handle the case where key already
			// exists in this map. If it does exist, we need to merge the vals
			//
			// update context as other blocks or attrs can refer to values we just
			// parsed
			ctx.Variables[key] = cty.ObjectVal(children)
		}
	}
}

// TODO(ryan-ph): This currently only handles locals and variable blocks. Other
// blocks will require defining / importing structs with struct tags for the
// necessary blocks.
func (h *HCLEditor) UpdateCtx(ctx *hcl.EvalContext) error {
	var r root
	if diags := gohcl.DecodeBody(h.hclFile.Body, ctx, &r); diags.HasErrors() {
		return diags
	}

	// TODO(ryan-ph): really this needs to do a topo-sort of _all_ of the blocks
	// held in Root to determine the order in which things should be evaluated.
	// This current implementation will not work well with blocks which refer
	// to other blocks that are defined later. For example, this will produce
	// an unknown value as we are currently just evaluating variables, then
	// locals top-to-bottom:
	// ```hcl
	// locals {
	//   team_name = "team-${local.owner}"
	// }
	// locals {
	//   owner = "admin"
	// }
	// ```
	h.updateVariables(ctx, r.Variables...)
	h.updateLocals(ctx, r.Locals...)
	return nil
}

// Create creates attributes and blocks matched with the given query
// with the given value. The value can be any type and it's transformed
// into HCL type inside.
func (h *HCLEditor) Create(queryStr string, value interface{}, opts ...Option) error {
	defer h.reload()

	opt := &option{}
	for _, optFunc := range opts {
		optFunc(opt)
	}

	queries, err := query.Build(queryStr)
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
// TODO(ryan-ph): this needs to take *hcl.EvalContext)
func (h *HCLEditor) Read(queryStr string, opts ...Option) (map[string]interface{}, error) {
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
	// TODO(ryan-ph): this needs to take *hcl.EvalContext)
	hdlr, err := handler.NewReadHandler(results)
	if err != nil {
		return nil, err
	}

	w := &walker.Walker{
		Handler: hdlr,
		Mode:    walker.Read,
	}

	if err := w.Walk(h.writeFile.Body(), queries, 0, []string{}); err != nil {
		return nil, err
	}

	return convert(results)
}

// Update replaces attributes and blocks which matched with its key
// and given query with the given value. The value can be any type
// and it's transformed into HCL type inside.
//
// By default, it returns error if the does not matched with any key.
// You must create value before update.
func (h *HCLEditor) Update(queryStr string, value interface{}, opts ...Option) error {
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
	} else if opt.beforeNewline {
		return fmt.Errorf("WithNewLine is not supported for Update")
	}

	queries, err := query.Build(queryStr)
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
