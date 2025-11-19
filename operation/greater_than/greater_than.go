package greater_than

import (
	"context"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/operation"
	"github.com/liyanbing/filter/utils"
	"github.com/liyanbing/filter/variable"
)

func init() {
	operation.Register(&GreaterThan{name: ">"})
	operation.Register(&GreaterThan{name: "gt"})
}

type GreaterThan struct {
	operation.BaseOperationPrepareValue
	name string
}

func (s *GreaterThan) Name() string { return s.name }
func (s *GreaterThan) Run(ctx context.Context, v variable.Variable, value interface{}, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variable.GetValue(ctx, v, data, cache)
	if err != nil {
		return false, err
	}
	return utils.ObjectCompare(variableValue, value) == 1, nil
}
