package has

import (
	"context"
	"fmt"

	"github.com/liyanbing/filter/cache"
	filterType "github.com/liyanbing/filter/filter_type"
	"github.com/liyanbing/filter/operation"
	"github.com/liyanbing/filter/utils"
	"github.com/liyanbing/filter/variable"
)

const Name = "has"

func init() {
	operation.Register(&Has{})
}

type Has struct{}

func (s *Has) Name() string { return Name }
func (s *Has) PrepareValue(value interface{}) (interface{}, error) {
	targetValues := utils.ParseTargetArrayValue(value)
	if len(targetValues) == 0 {
		return nil, fmt.Errorf("[%s] expression value must be greater than one element", Name)
	}
	return targetValues, nil
}

func (s *Has) Run(ctx context.Context, v variable.Variable, value interface{}, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variable.GetValue(ctx, v, data, cache)
	if err != nil {
		return false, nil
	}

	targetValueElements, ok := value.([]interface{})
	if !ok {
		return false, fmt.Errorf("[%s] value must be an array", Name)
	}

	variableValueElements := utils.ParseTargetArrayValue(variableValue)
	for _, valueElement := range targetValueElements {
		has := false
		for _, variableValueElement := range variableValueElements {
			if filterType.ObjectCompare(valueElement, variableValueElement) == 0 {
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
