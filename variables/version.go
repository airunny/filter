package variables

import (
	"context"
	"errors"

	"github.com/liyanbing/filter/cache"
	filterContext "github.com/liyanbing/filter/context"
)

const versionName = "version"

func versionVariable() *Version {
	return &Version{}
}

// Version 应用版本
type Version struct{ CacheableVariable }

func (s *Version) Name() string { return versionName }

func (s *Version) Value(ctx context.Context, _ interface{}, _ *cache.Cache) (interface{}, error) {
	version, ok := filterContext.FromVersion(ctx)
	if !ok {
		return nil, errors.New("version not found in context")
	}
	return version, nil
}
