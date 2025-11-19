package between

import (
	"context"
	"fmt"

	"github.com/liyanbing/filter/cache"
	filterType "github.com/liyanbing/filter/filter_type"
	"github.com/liyanbing/filter/operation"
	"github.com/liyanbing/filter/utils"
	"github.com/liyanbing/filter/variable"
)

const Name = "between"

func init() {
	operation.Register(&Between{})
}

type Between struct{}

func (s *Between) Name() string { return s.Name() }
func (s *Between) PrepareValue(value interface{}) (interface{}, error) {
	elements := utils.ParseTargetArrayValue(value)
	if len(elements) != 2 {
		return nil, fmt.Errorf("[%s] expression must have two element", Name)
	}
	return elements, nil
}

func (s *Between) Run(ctx context.Context, v variable.Variable, value interface{}, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variable.GetValue(ctx, v, data, cache)
	if err != nil {
		return false, err
	}
	startAndEnd := value.([]interface{})
	return filterType.ObjectCompare(variableValue, startAndEnd[0]) >= 0 && filterType.ObjectCompare(variableValue, startAndEnd[1]) <= 0, nil
}
