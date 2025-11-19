package not

import (
	"context"
	"fmt"

	"github.com/liyanbing/filter/cache"
	filterType "github.com/liyanbing/filter/filter_type"
	"github.com/liyanbing/filter/operation"
	"github.com/liyanbing/filter/utils"
	"github.com/liyanbing/filter/variable"
)

const Name = "not"

func init() {
	operation.Register(&Not{})
}

type Not struct{}

func (s *Not) Name() string { return Name }
func (s *Not) PrepareValue(value interface{}) (interface{}, error) {
	targetValues := utils.ParseTargetArrayValue(value)
	if len(targetValues) == 0 {
		return nil, fmt.Errorf("[%s] operation value must be greater than one element", Name)
	}
	return targetValues, nil
}

func (s *Not) Run(ctx context.Context, v variable.Variable, value interface{}, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variable.GetValue(ctx, v, data, cache)
	if err != nil {
		return false, nil
	}

	variableValueElements := utils.ParseTargetArrayValue(variableValue)
	targetValueElements, ok := value.([]interface{})
	if !ok {
		return false, fmt.Errorf("[%s] value must be array", Name)
	}

	for _, variableValueElement := range variableValueElements {
		for _, valueElement := range targetValueElements {
			if filterType.ObjectCompare(variableValueElement, valueElement) == 0 {
				return false, nil
			}
		}
	}
	return true, nil
}
