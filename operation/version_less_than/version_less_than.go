package version_less_than

import (
	"context"

	"github.com/liyanbing/filter/cache"
	filterType "github.com/liyanbing/filter/filter_type"
	"github.com/liyanbing/filter/operation"
	"github.com/liyanbing/filter/variable"
	"github.com/liyanbing/filter/version"
)

const Name = "vlt"

func init() {
	operation.Register(&VersionLessThan{})
}

type VersionLessThan struct {
	operation.BaseOperationPrepareValue
}

func (s *VersionLessThan) Name() string { return Name }
func (s *VersionLessThan) Run(ctx context.Context, v variable.Variable, value interface{}, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variable.GetValue(ctx, v, data, cache)
	if err != nil {
		return false, err
	}

	if version.Compare(filterType.GetString(variableValue), filterType.GetString(value)) < 0 {
		return true, nil
	}
	return false, nil
}
