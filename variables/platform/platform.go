package variables

import (
	"context"
	"errors"

	"github.com/liyanbing/filter/cache"
	filterContext "github.com/liyanbing/filter/context"
)

const platformName = "platform"

func platformVariable() *Platform {
	return &Platform{}
}

// Platform 平台
type Platform struct{ CacheableVariable }

func (s *Platform) Name() string { return platformName }
func (s *Platform) Value(ctx context.Context, _ interface{}, _ *cache.Cache) (interface{}, error) {
	plt, ok := filterContext.FromPlatform(ctx)
	if !ok {
		return nil, errors.New("platform not found in context")
	}
	return plt, nil
}
