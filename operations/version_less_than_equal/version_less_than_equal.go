package version_less_than_equal

import (
	"context"
	"go/version"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/operations"
	"github.com/liyanbing/filter/types"
	"github.com/liyanbing/filter/variables"
)

const Name = "vlte"

func init() {
	operations.Register(&VersionLessThanEqual{})
}

type VersionLessThanEqual struct {
	operations.OriginValue
}

func (s *VersionLessThanEqual) Name() string { return Name }
func (s *VersionLessThanEqual) Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variables.GetValue(ctx, variable, data, cache)
	if err != nil {
		return false, err
	}
	if version.Compare(types.GetString(variableValue), types.GetString(value)) <= 0 {
		return true, nil
	}
	return false, nil
}
