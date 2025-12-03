package less_than_equal

import (
	"context"

	"github.com/airunny/filter/cache"
	"github.com/airunny/filter/operations"
	"github.com/airunny/filter/utils"
	"github.com/airunny/filter/variables"
)

func init() {
	operations.Register(&LessThanEqual{name: "<="})
	operations.Register(&LessThanEqual{name: "lte"})
}

type LessThanEqual struct {
	operations.OriginValue
	name string
}

func (s *LessThanEqual) Name() string { return s.name }
func (s *LessThanEqual) Run(ctx context.Context, variable variables.Variable, operationValue, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variables.GetValue(ctx, variable, data, cache)
	if err != nil {
		return false, err
	}
	return utils.ObjectCompare(variableValue, operationValue) <= 0, nil
}
