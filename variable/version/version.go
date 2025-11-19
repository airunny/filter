package version

import (
	"context"
	"errors"

	"github.com/liyanbing/filter/cache"
	filterContext "github.com/liyanbing/filter/context"
	"github.com/liyanbing/filter/variable"
)

const Name = "version"

func init() {
	variable.Register(variable.NewSimpleVariable(&Version{}))
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
