package not_equal

import (
	"context"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/operation"
	"github.com/liyanbing/filter/utils"
	"github.com/liyanbing/filter/variable"
)

func init() {
	operation.Register(&NotEqual{name: "!="})
	operation.Register(&NotEqual{name: "<>"})
	operation.Register(&NotEqual{name: "ne"})
}

type NotEqual struct {
	operation.BaseOperationPrepareValue
	name string
}

func (s *NotEqual) Name() string { return s.name }
func (s *NotEqual) Run(ctx context.Context, v variable.Variable, value interface{}, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variable.GetValue(ctx, v, data, cache)
	if err != nil {
		return false, err
	}
	return utils.ObjectCompare(variableValue, value) != 0, nil
}
