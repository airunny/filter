package version_less_than

import (
	"context"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/operations"
	"github.com/liyanbing/filter/types"
	"github.com/liyanbing/filter/utils"
	"github.com/liyanbing/filter/variables"
)

const Name = "vlt"

func init() {
	operations.Register(&VersionLessThan{})
}

type VersionLessThan struct {
	operations.OriginValue
}

func (s *VersionLessThan) Name() string { return Name }
func (s *VersionLessThan) Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variables.GetValue(ctx, variable, data, cache)
	if err != nil {
		return false, err
	}

	if utils.VersionCompare(types.GetString(variableValue), types.GetString(value)) < 0 {
		return true, nil
	}
	return false, nil
}
