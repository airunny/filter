package version_greater_than

import (
	"context"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/operations"
	"github.com/liyanbing/filter/types"
	"github.com/liyanbing/filter/utils"
	"github.com/liyanbing/filter/variables"
)

const Name = "vgt"

func init() {
	operations.Register(&VersionGreaterThan{})
}

type VersionGreaterThan struct {
	operations.OriginValue
}

func (s *VersionGreaterThan) Name() string { return Name }
func (s *VersionGreaterThan) Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variables.GetValue(ctx, variable, data, cache)
	if err != nil {
		return false, err
	}

	if utils.VersionCompare(types.GetString(variableValue), types.GetString(value)) == 1 {
		return true, nil
	}
	return false, nil
}
