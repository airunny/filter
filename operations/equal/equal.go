package equal

import (
	"context"

	"github.com/airunny/filter/cache"
	"github.com/airunny/filter/operations"
	"github.com/airunny/filter/utils"
	"github.com/airunny/filter/variables"
)

func init() {
	operations.Register(&Equal{name: "="})
	operations.Register(&Equal{name: "eq"})
}

type Equal struct {
	operations.OriginValue
	name string
}

func (s *Equal) Name() string { return s.name }
func (s *Equal) Run(ctx context.Context, variable variables.Variable, operationValue, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variables.GetValue(ctx, variable, data, cache)
	if err != nil {
		return false, err
	}
	return utils.ObjectCompare(variableValue, operationValue) == 0, nil
}
