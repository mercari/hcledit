package walker

import (
	"github.com/hashicorp/hcl/v2/hclwrite"

	"go.mercari.io/hcledit/internal/ast"
	"go.mercari.io/hcledit/internal/handler"
	"go.mercari.io/hcledit/internal/query"
)

type walkMode int

const (
	Create walkMode = iota
	Read
	Update
	Delete
)

type Walker struct {
	Handler handler.Handler
	Mode    walkMode
}

func (w *Walker) Walk(body *hclwrite.Body, queries []query.Query, index int, keytrail []string) error {
	if err := w.walkAttribute(body, queries, index, keytrail); err != nil {
		return err
	}

	if err := w.walkBlock(body, queries, index, keytrail); err != nil {
		return err
	}

	return nil
}

func (w *Walker) walkAttribute(body *hclwrite.Body, queries []query.Query, index int, keytrail []string) error {
	var handled bool
	for key := range body.Attributes() {
		if queries[index].Match(key) {
			keytrail = append(keytrail, key)

			// This means it reaches to the end of queries where
			// we should execute the handler.
			if index == len(queries)-1 {
				switch w.Mode {
				case Delete:
					body.RemoveAttribute(key)
				default:
					if err := w.Handler.HandleBody(body, key, keytrail); err != nil {
						return err
					}
				}
				handled = true
			} else {
				// This means the query indicates more room to go deeper.
				//
				// It tries to parse it as Object type and if it works
				// then walk into the Object.
				nestedIndex := index
				nestedIndex++

				obj, err := ast.ParseObject(body.GetAttribute(key).BuildTokens(nil))
				if err != nil {
					continue
				}

				handled, err = w.walkObject(obj, queries, nestedIndex, keytrail)
				if err != nil {
					return err
				}
				if !handled {
					continue
				}

				// If object is modified while working, we need to add it as new attribute.
				// Since `hclwrite` provides only `.AppendUnstructuredTokens` function,
				// to preserve the original order of attribute, we need to extract it
				// and do special operation to insert.
				tokens := obj.BuildTokens()
				ast.ReplaceBodyTokens(body, key, tokens)
			}
		}
	}

	// Create a new attribute when it meets the following conditions:
	//   - It reaches to the end of queries
	//   - It does not do any operation
	//   - It's executed as creation mode
	if index == len(queries)-1 && !handled && w.Mode == Create {
		if key := queries[index].Key(); key != "*" {
			keytrail = append(keytrail, key)
			if err := w.Handler.HandleBody(body, key, keytrail); err != nil {
				return err
			}
		}
	}

	return nil
}

func (w *Walker) walkBlock(body *hclwrite.Body, queries []query.Query, index int, keytrail []string) error {
	for _, block := range body.Blocks() {
		blockIndex := index
		if !queries[blockIndex].Match(block.Type()) {
			continue
		}

		blockIndex++
		keytrail := append(keytrail, block.Type())

		var unmatched bool
		if blockIndex <= len(queries)-1 {
			for _, label := range block.Labels() {
				if queries[blockIndex].Match(label) {
					blockIndex++
					keytrail = append(keytrail, label)
					continue
				}
				unmatched = true
			}
		}

		if unmatched {
			continue
		}

		if blockIndex <= len(queries)-1 {
			// This means the query indicates more room to go deeper.
			if err := w.Walk(block.Body(), queries, blockIndex, keytrail); err != nil {
				return err
			}
		} else if index == len(queries)-1 {
			// This means it reaches to the end of queries where
			// we should execute the handler.
			switch w.Mode {
			case Delete:
				body.RemoveBlock(block)
			case Update:
				if err := w.Handler.HandleBody(body, "", keytrail); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (w *Walker) walkObject(obj *ast.Object, queries []query.Query, index int, keytrail []string) (bool, error) {
	var handled bool
	for key := range obj.ObjectAttributes() {
		if queries[index].Match(key) {
			keytrail = append(keytrail, key)

			// This means it reaches to the end of queries where
			// we should execute the handler.
			if index == len(queries)-1 {
				switch w.Mode {
				case Delete:
					obj.DeleteObjectAttribute(key)
				default:
					if err := w.Handler.HandleObject(obj, key, keytrail); err != nil {
						return false, err
					}
				}
				handled = true
			} else {
				// This means the query indicates more room to go deeper and
				// object has nested object.
				//
				// It tries to parse it as Object type and if it works
				// then walk into the Object.
				nestedIndex := index
				nestedIndex++

				nestedObj, err := ast.ParseObject(obj.GetObjectAttribute(key).BuildTokens())
				if err != nil {
					continue
				}

				handled, err = w.walkObject(nestedObj, queries, nestedIndex, keytrail)
				if err != nil {
					return false, err
				}
				if !handled {
					continue
				}

				// If attribute of nested object is updated, we need to replace its object
				obj.ReplaceObjectAttribute(key, nestedObj)
			}
		}
	}

	// Create a new attribute when it meets the following conditions:
	//   - It reaches to the end of queries
	//   - It does not do any operation
	//   - It's executed as creation mode
	if index == len(queries)-1 && !handled && w.Mode == Create {
		if key := queries[index].Key(); key != "*" {
			keytrail = append(keytrail, key)
			if err := w.Handler.HandleObject(obj, key, keytrail); err != nil {
				return false, err
			}
			handled = true
		}
	}

	return handled, nil
}
