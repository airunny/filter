package has

import (
	"context"
	"fmt"

	"github.com/airunny/filter/cache"
	"github.com/airunny/filter/operations"
	"github.com/airunny/filter/utils"
	"github.com/airunny/filter/variables"
)

const Name = "has"

var ErrInvalidOperationValue = fmt.Errorf("[%s] operation value must be greater than one element", Name)

func init() {
	operations.Register(&Has{})
}

type Has struct{}

func (s *Has) Name() string { return Name }
func (s *Has) PrepareValue(value interface{}) (interface{}, error) {
	targetValues := utils.ParseTargetArrayValue(value)
	if len(targetValues) == 0 {
		return nil, ErrInvalidOperationValue
	}
	return targetValues, nil
}

func (s *Has) Run(ctx context.Context, variable variables.Variable, operationValue, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variables.GetValue(ctx, variable, data, cache)
	if err != nil {
		return false, err
	}

	targetValueElements, ok := operationValue.([]interface{})
	if !ok {
		return false, ErrInvalidOperationValue
	}

	variableValueElements := utils.ParseTargetArrayValue(variableValue)
	for _, valueElement := range targetValueElements {
		has := false
		for _, variableValueElement := range variableValueElements {
			if utils.ObjectCompare(valueElement, variableValueElement) == 0 {
				has = true
				break
			}
		}
		if !has {
			return false, nil
		}
	}
	return true, nil
}
