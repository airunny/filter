package not_in

import (
	"context"
	"fmt"

	"github.com/airunny/filter/cache"
	"github.com/airunny/filter/operations"
	"github.com/airunny/filter/utils"
	"github.com/airunny/filter/variables"
)

const Name = "nin"

var emptyElementErr = fmt.Errorf("[%s] operation value must be greater than one element", Name)

func init() {
	operations.Register(&NotIn{})
}

type NotIn struct{}

func (s *NotIn) Name() string { return Name }
func (s *NotIn) PrepareValue(value interface{}) (interface{}, error) {
	targetValues := utils.ParseTargetArrayValue(value)
	if len(targetValues) == 0 {
		return nil, emptyElementErr
	}
	return targetValues, nil
}
func (s *NotIn) Run(ctx context.Context, variable variables.Variable, operationValue, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variables.GetValue(ctx, variable, data, cache)
	if err != nil {
		return false, err
	}

	if targetValues, ok := operationValue.([]interface{}); ok {
		for _, targetValue := range targetValues {
			if utils.ObjectCompare(variableValue, targetValue) == 0 {
				return false, nil
			}
		}
	}
	return true, nil
}
