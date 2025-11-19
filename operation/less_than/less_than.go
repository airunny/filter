package less_than

import (
	"context"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/operation"
	"github.com/liyanbing/filter/utils"
	"github.com/liyanbing/filter/variable"
)

func init() {
	operation.Register(&LessThan{name: "<"})
	operation.Register(&LessThan{name: "lt"})
}

type LessThan struct {
	operation.BaseOperationPrepareValue
	name string
}

func (s *LessThan) Name() string { return s.name }
func (s *LessThan) Run(ctx context.Context, v variable.Variable, value interface{}, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variable.GetValue(ctx, v, data, cache)
	if err != nil {
		return false, err
	}
	return utils.ObjectCompare(variableValue, value) == -1, nil
}
