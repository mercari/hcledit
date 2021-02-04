package converter

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
	"github.com/zclconf/go-cty/cty/gocty"
)

func FromCtyValueToGoValue(ctyVal cty.Value) (interface{}, error) {
	switch ctyVal.Type() {
	case cty.Number:
		var goVal int
		err := gocty.FromCtyValue(ctyVal, &goVal)
		return goVal, err
	case cty.String:
		var goVal string
		err := gocty.FromCtyValue(ctyVal, &goVal)
		return goVal, err
	case cty.Bool:
		var goVal bool
		err := gocty.FromCtyValue(ctyVal, &goVal)
		return goVal, err
	}

	{
		var goVal []int
		convTy, err := gocty.ImpliedType(goVal)
		if err != nil {
			panic(fmt.Sprintf("Should not be reached here: %s", err))
		}

		srcVal, err := convert.Convert(ctyVal, convTy)
		if err == nil {
			if err := gocty.FromCtyValue(srcVal, &goVal); err == nil {
				return goVal, nil
			}
		}
	}

	{
		var goVal []bool
		convTy, err := gocty.ImpliedType(goVal)
		if err != nil {
			panic(fmt.Sprintf("Should not be reached here: %s", err))
		}

		srcVal, err := convert.Convert(ctyVal, convTy)
		if err == nil {
			if err := gocty.FromCtyValue(srcVal, &goVal); err == nil {
				return goVal, nil
			}
		}
	}

	{
		var goVal []string
		convTy, err := gocty.ImpliedType(goVal)
		if err != nil {
			panic(fmt.Sprintf("Should not be reached here: %s", err))
		}

		srcVal, err := convert.Convert(ctyVal, convTy)
		if err == nil {
			if err := gocty.FromCtyValue(srcVal, &goVal); err == nil {
				return goVal, nil
			}
		}
	}

	return nil, fmt.Errorf("unsupported cty type")
}
