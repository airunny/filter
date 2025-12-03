package version

import (
	"context"
	"errors"

	"github.com/airunny/filter/cache"
	filterContext "github.com/airunny/filter/context"
	"github.com/airunny/filter/variables"
)

const Name = "version"

func init() {
	variables.Register(variables.NewSimpleVariable(&Version{}))
}

// Version 应用版本
type Version struct{}

func (s *Version) Name() string    { return Name }
func (s *Version) Cacheable() bool { return true }
func (s *Version) Value(ctx context.Context, _ interface{}, _ *cache.Cache) (interface{}, error) {
	version, ok := filterContext.FromVersion(ctx)
	if !ok {
		return nil, errors.New("version not found in context")
	}
	return version, nil
}
