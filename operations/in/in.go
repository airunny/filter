package in

import (
	"context"
	"fmt"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/operations"
	"github.com/liyanbing/filter/utils"
	"github.com/liyanbing/filter/variables"
)

const Name = "in"

var emptyElementErr = fmt.Errorf("[%s] operation value must be greater than one element", Name)

func init() {
	operations.Register(&In{})
}

type In struct{}

func (s *In) Name() string { return Name }
func (s *In) PrepareValue(value interface{}) (interface{}, error) {
	targetValues := utils.ParseTargetArrayValue(value)
	if len(targetValues) == 0 {
		return nil, emptyElementErr
	}
	return targetValues, nil
}

func (s *In) Run(ctx context.Context, variable variables.Variable, operationValue, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variables.GetValue(ctx, variable, data, cache)
	if err != nil {
		return false, err
	}

	variableValueElements := utils.ParseTargetArrayValue(variableValue)
	if targetValues, ok := operationValue.([]interface{}); ok {
		for _, variableValueElement := range variableValueElements {
			exists := false
			for _, targetValue := range targetValues {
				if utils.ObjectCompare(variableValueElement, targetValue) == 0 {
					exists = true
				}
			}

			if !exists {
				return false, nil
			}
		}
		return true, nil
	}
	return false, nil
}
