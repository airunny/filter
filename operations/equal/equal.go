package equal

import (
	"context"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/operations"
	"github.com/liyanbing/filter/utils"
	"github.com/liyanbing/filter/variables"
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
func (s *Equal) Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variables.GetValue(ctx, variable, data, cache)
	if err != nil {
		return false, err
	}
	return utils.ObjectCompare(variableValue, value) == 0, nil
}
