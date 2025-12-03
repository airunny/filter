package any

import (
	"context"
	"fmt"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/operations"
	"github.com/liyanbing/filter/utils"
	"github.com/liyanbing/filter/variables"
)

const Name = "any"

var ErrInvalidOperationValue = fmt.Errorf("[%s] operation value must be greater than one element", Name)

func init() {
	operations.Register(&Any{})
}

type Any struct{}

func (s *Any) Name() string { return Name }
func (s *Any) PrepareValue(value interface{}) (interface{}, error) {
	targetValues := utils.ParseTargetArrayValue(value)
	if len(targetValues) == 0 {
		return nil, ErrInvalidOperationValue
	}
	return targetValues, nil
}

func (s *Any) Run(ctx context.Context, variable variables.Variable, operationValue, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variables.GetValue(ctx, variable, data, cache)
	if err != nil {
		return false, err
	}

	variableValueElements := utils.ParseTargetArrayValue(variableValue)
	targetValueElements, ok := operationValue.([]interface{})
	if !ok {
		return false, ErrInvalidOperationValue
	}

	for _, targetValueElement := range targetValueElements {
		for _, variableValueElement := range variableValueElements {
			if utils.ObjectCompare(targetValueElement, variableValueElement) == 0 {
				return true, nil
			}
		}
	}
	return false, nil
}
