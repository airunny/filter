package version_less_than_equal

import (
	"context"
	"go/version"

	"github.com/liyanbing/filter/cache"
	filterType "github.com/liyanbing/filter/filter_type"
	"github.com/liyanbing/filter/operation"
	"github.com/liyanbing/filter/variable"
)

const Name = "vlte"

func init() {
	operation.Register(&VersionLessThanEqual{})
}

type VersionLessThanEqual struct {
	operation.BaseOperationPrepareValue
}

func (s *VersionLessThanEqual) Name() string { return Name }
func (s *VersionLessThanEqual) Run(ctx context.Context, v variable.Variable, value interface{}, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variable.GetValue(ctx, v, data, cache)
	if err != nil {
		return false, err
	}
	if version.Compare(filterType.GetString(variableValue), filterType.GetString(value)) <= 0 {
		return true, nil
	}
	return false, nil
}
