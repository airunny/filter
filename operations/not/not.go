package not

import (
	"context"
	"fmt"

	"github.com/airunny/filter/cache"
	"github.com/airunny/filter/operations"
	"github.com/airunny/filter/utils"
	"github.com/airunny/filter/variables"
)

const Name = "not"

var ErrInvalidOperationValue = fmt.Errorf("[%s] operation value must be greater than one element", Name)

func init() {
	operations.Register(&Not{})
}

type Not struct{}

func (s *Not) Name() string { return Name }
func (s *Not) PrepareValue(value interface{}) (interface{}, error) {
	targetValues := utils.ParseTargetArrayValue(value)
	if len(targetValues) == 0 {
		return nil, ErrInvalidOperationValue
	}
	return targetValues, nil
}

func (s *Not) Run(ctx context.Context, variable variables.Variable, operationValue, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variables.GetValue(ctx, variable, data, cache)
	if err != nil {
		return false, err
	}

	variableValueElements := utils.ParseTargetArrayValue(variableValue)
	targetValueElements, ok := operationValue.([]interface{})
	if !ok {
		return false, ErrInvalidOperationValue
	}

	for _, variableValueElement := range variableValueElements {
		for _, valueElement := range targetValueElements {
			if utils.ObjectCompare(variableValueElement, valueElement) == 0 {
				return false, nil
			}
		}
	}
	return true, nil
}
