package any

import (
	"context"
	"fmt"

	"github.com/liyanbing/filter/cache"
	filterType "github.com/liyanbing/filter/filter_type"
	"github.com/liyanbing/filter/operation"
	"github.com/liyanbing/filter/utils"
	"github.com/liyanbing/filter/variable"
)

const Name = "any"

func init() {
	operation.Register(&Any{})
}

type Any struct{}

func (s *Any) Name() string { return Name }
func (s *Any) PrepareValue(value interface{}) (interface{}, error) {
	targetValues := utils.ParseTargetArrayValue(value)
	if len(targetValues) == 0 {
		return nil, fmt.Errorf("[%s] expression must be greater than one element", Name)
	}
	return targetValues, nil
}

func (s *Any) Run(ctx context.Context, v variable.Variable, value interface{}, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variable.GetValue(ctx, v, data, cache)
	if err != nil {
		return false, err
	}

	variableValueElements := utils.ParseTargetArrayValue(variableValue)
	targetValueElements, ok := value.([]interface{})
	if !ok {
		return false, fmt.Errorf("[%s] value must be an array", Name)
	}

	for _, targetValueElement := range targetValueElements {
		for _, variableValueElement := range variableValueElements {
			if filterType.ObjectCompare(targetValueElement, variableValueElement) == 0 {
				return true, nil
			}
		}
	}
	return false, nil
}
