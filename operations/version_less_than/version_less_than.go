package version_less_than

import (
	"context"

	"github.com/airunny/filter/cache"
	"github.com/airunny/filter/operations"
	"github.com/airunny/filter/types"
	"github.com/airunny/filter/utils"
	"github.com/airunny/filter/variables"
)

const Name = "vlt"

func init() {
	operations.Register(&VersionLessThan{})
}

type VersionLessThan struct {
	operations.OriginValue
}

func (s *VersionLessThan) Name() string { return Name }
func (s *VersionLessThan) Run(ctx context.Context, variable variables.Variable, operationValue, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variables.GetValue(ctx, variable, data, cache)
	if err != nil {
		return false, err
	}

	if utils.VersionCompare(types.GetString(variableValue), types.GetString(operationValue)) < 0 {
		return true, nil
	}
	return false, nil
}
