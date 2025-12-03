package not_equal

import (
	"context"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/operations"
	"github.com/liyanbing/filter/utils"
	"github.com/liyanbing/filter/variables"
)

func init() {
	operations.Register(&NotEqual{name: "!="})
	operations.Register(&NotEqual{name: "<>"})
	operations.Register(&NotEqual{name: "ne"})
}

type NotEqual struct {
	operations.OriginValue
	name string
}

func (s *NotEqual) Name() string { return s.name }
func (s *NotEqual) Run(ctx context.Context, variable variables.Variable, operationValue, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variables.GetValue(ctx, variable, data, cache)
	if err != nil {
		return false, err
	}
	return utils.ObjectCompare(variableValue, operationValue) != 0, nil
}
