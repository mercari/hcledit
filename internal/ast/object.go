package ast

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Object
type Object struct {
	beforeTokens     hclwrite.Tokens
	objectAtrributes []*ObjectAtrribute
	afterTokens      hclwrite.Tokens
}

func (o *Object) BuildTokens() hclwrite.Tokens {
	var tokens hclwrite.Tokens
	tokens = append(tokens, o.beforeTokens...)
	for _, oa := range o.objectAtrributes {
		tokens = append(tokens, oa.BuildTokens()...)
	}
	tokens = append(tokens, o.afterTokens...)
	return tokens
}

func (o *Object) ObjectAttributes() map[string]*ObjectAtrribute {
	result := make(map[string]*ObjectAtrribute)
	for _, objAttr := range o.objectAtrributes {
		result[objAttr.name] = objAttr
	}
	return result
}

func (o *Object) objectAttributeKeys() []string {
	keys := make([]string, 0, len(o.objectAtrributes))
	for _, objAttr := range o.objectAtrributes {
		keys = append(keys, objAttr.name)
	}
	return keys
}

func (o *Object) GetObjectAttribute(name string) *ObjectAtrribute {
	for _, objAttr := range o.objectAtrributes {
		if objAttr.name == name {
			return objAttr
		}
	}
	return nil
}

func (o *Object) SetObjectAttributeRaw(name string, exprTokens, beforeTokens hclwrite.Tokens) *ObjectAtrribute {
	objAttr := o.GetObjectAttribute(name)
	if objAttr != nil {
		objAttr.exprTokens = exprTokens
	} else {
		objAttr = newObjectAttribute(name, exprTokens, beforeTokens)
		o.objectAtrributes = append(o.objectAtrributes, objAttr)
	}
	return objAttr
}

func (o *Object) DeleteObjectAttribute(name string) *ObjectAtrribute {
	var (
		targetIdx     int
		targetObjAttr *ObjectAtrribute
	)

	for idx, objAttr := range o.objectAtrributes {
		if objAttr.name == name {
			targetIdx = idx
			targetObjAttr = objAttr
		}
	}

	// Delete item from slice
	o.objectAtrributes = append(o.objectAtrributes[:targetIdx], o.objectAtrributes[targetIdx+1:]...)

	return targetObjAttr
}

func (o *Object) ReplaceObjectAttribute(name string, nestedObject *Object) *ObjectAtrribute {
	objAttr := o.GetObjectAttribute(name)
	if objAttr != nil {
		objAttr.beforeTokens = nestedObject.beforeTokens
		objAttr.exprTokens = nestedObject.expr()
		objAttr.afterTokens = nestedObject.afterTokens
	}
	return objAttr
}

func (o *Object) UpdateObjectAttributeOrder(targetKey, afterKey string) {
	objectAtrributes := make([]*ObjectAtrribute, 0, len(o.objectAtrributes))

	var matched bool
	for _, originalKey := range o.objectAttributeKeys() {
		if originalKey == targetKey {
			continue
		}
		objectAtrributes = append(objectAtrributes, o.GetObjectAttribute(originalKey))
		if originalKey == afterKey {
			objectAtrributes = append(objectAtrributes, o.GetObjectAttribute(targetKey))
			matched = true
		}
	}

	if matched {
		o.objectAtrributes = objectAtrributes
	}
}

func (o *Object) expr() hclwrite.Tokens {
	var tokens hclwrite.Tokens
	for _, oa := range o.objectAtrributes {
		tokens = append(tokens, oa.BuildTokens()...)
	}
	return tokens
}
