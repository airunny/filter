package variables

import (
	"context"
	"errors"

	"github.com/liyanbing/filter/cache"
	filterContext "github.com/liyanbing/filter/context"
)

const uaName = "ua"

func uaVariable() *UserAgent {
	return &UserAgent{}
}

// UserAgent 用户代理信息
type UserAgent struct{ CacheableVariable }

func (s *UserAgent) Name() string { return uaName }
func (s *UserAgent) Value(ctx context.Context, _ interface{}, _ *cache.Cache) (interface{}, error) {
	ua, ok := filterContext.FromUA(ctx)
	if !ok {
		return nil, errors.New("ua not found in context")
	}
	return ua, nil
}
