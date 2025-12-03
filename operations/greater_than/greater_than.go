package greater_than

import (
	"context"

	"github.com/airunny/filter/cache"
	"github.com/airunny/filter/operations"
	"github.com/airunny/filter/utils"
	"github.com/airunny/filter/variables"
)

func init() {
	operations.Register(&GreaterThan{name: ">"})
	operations.Register(&GreaterThan{name: "gt"})
}

type GreaterThan struct {
	operations.OriginValue
	name string
}

func (s *GreaterThan) Name() string { return s.name }
func (s *GreaterThan) Run(ctx context.Context, variable variables.Variable, operationValue, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variables.GetValue(ctx, variable, data, cache)
	if err != nil {
		return false, err
	}
	return utils.ObjectCompare(variableValue, operationValue) == 1, nil
}
