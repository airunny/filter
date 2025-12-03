package less_than

import (
	"context"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/operations"
	"github.com/liyanbing/filter/utils"
	"github.com/liyanbing/filter/variables"
)

func init() {
	operations.Register(&LessThan{name: "<"})
	operations.Register(&LessThan{name: "lt"})
}

type LessThan struct {
	operations.OriginValue
	name string
}

func (s *LessThan) Name() string { return s.name }
func (s *LessThan) Run(ctx context.Context, variable variables.Variable, operationValue, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variables.GetValue(ctx, variable, data, cache)
	if err != nil {
		return false, err
	}
	return utils.ObjectCompare(variableValue, operationValue) == -1, nil
}
