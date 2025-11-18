package variables

import (
	"context"

	"github.com/liyanbing/filter/cache"
	filterContext "github.com/liyanbing/filter/context"
)

const isLoginName = "is_login"

func isLoginVariable() *IsLogin {
	return &IsLogin{}
}

// IsLogin 是否登录
type IsLogin struct{ CacheableVariable }

func (s *IsLogin) Name() string { return isLoginName }
func (s *IsLogin) Value(ctx context.Context, _ interface{}, _ *cache.Cache) (interface{}, error) {
	uid, ok := filterContext.FromUserId(ctx)
	if !ok {
		return false, nil
	}

	if uid != "" {
		return true, nil
	} else {
		return false, nil
	}
}
