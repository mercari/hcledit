package converter

import (
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

func FromCtyValueToGoValue(ctyVal cty.Value) (interface{}, error) {
	defer func() {
		if r := recover(); r != nil {
			// Return raw string representation if conversion fails
			ctyVal = cty.StringVal(ctyVal.GoString())
		}
	}()
	if ctyVal.Type() == cty.Number {
		if ctyVal.RawEquals(cty.NumberIntVal(0)) || ctyVal.RawEquals(cty.NumberIntVal(1)) {
			var goVal int
			err := gocty.FromCtyValue(ctyVal, &goVal)
			return goVal, err
		}
		var goVal float64
		err := gocty.FromCtyValue(ctyVal, &goVal)
		return goVal, err
	}
	switch ctyVal.Type() {
	case cty.String:
		var goVal string
		err := gocty.FromCtyValue(ctyVal, &goVal)
		return goVal, err
	case cty.Bool:
		var goVal bool
		err := gocty.FromCtyValue(ctyVal, &goVal)
		return goVal, err
	default:
		if ctyVal.Type().IsListType() || ctyVal.Type().IsTupleType() {
			return convertSpecificList(ctyVal)
		} else if ctyVal.Type().IsMapType() {
			return convertMap(ctyVal)
		} else if ctyVal.Type().IsObjectType() {
			return convertObject(ctyVal)
		}
		return ctyVal.GoString(), nil
	}
}

func convertSpecificList(ctyVal cty.Value) (interface{}, error) {
	if !ctyVal.IsKnown() || ctyVal.IsNull() {
		return nil, nil
	}
	var result []interface{}
	for _, elem := range ctyVal.AsValueSlice() {
		converted, err := FromCtyValueToGoValue(elem)
		if err != nil {
			return nil, err
		}
		result = append(result, converted)
	}

	// Determine specific type of list
	if len(result) > 0 {
		switch result[0].(type) {
		case bool:
			boolList := make([]bool, len(result))
			for i, v := range result {
				boolList[i] = v.(bool)
			}
			return boolList, nil
		case string:
			stringList := make([]string, len(result))
			for i, v := range result {
				stringList[i] = v.(string)
			}
			return stringList, nil
		case float64:
			floatList := make([]float64, len(result))
			for i, v := range result {
				floatList[i] = v.(float64)
			}
			return floatList, nil
		case int:
			intList := make([]int, len(result))
			for i, v := range result {
				if f, ok := v.(float64); ok {
					intList[i] = int(f)
				} else {
					intList[i] = v.(int)
				}
			}
			return intList, nil
		}
	}
	return result, nil
}

func convertList(ctyVal cty.Value) ([]interface{}, error) {
	if !ctyVal.IsKnown() || ctyVal.IsNull() {
		return nil, nil
	}
	var result []interface{}
	for _, elem := range ctyVal.AsValueSlice() {
		converted, err := FromCtyValueToGoValue(elem)
		if err != nil {
			return nil, err
		}
		result = append(result, converted)
	}
	return result, nil
}

func convertMap(ctyVal cty.Value) (map[string]interface{}, error) {
	if !ctyVal.IsKnown() || ctyVal.IsNull() {
		return nil, nil
	}
	result := make(map[string]interface{})
	for key, value := range ctyVal.AsValueMap() {
		converted, err := FromCtyValueToGoValue(value)
		if err != nil {
			return nil, err
		}
		result[key] = converted
	}
	return result, nil
}

func convertObject(ctyVal cty.Value) (map[string]interface{}, error) {
	return convertMap(ctyVal)
}
