package variables

import (
	"context"

	"github.com/liyanbing/filter/cache"
	filterContext "github.com/liyanbing/filter/context"
	"github.com/liyanbing/filter/variable"
)

const Name = "is_login"

func init() {
	variable.Register(variable.NewSimpleVariable(&IsLogin{}))
}

// IsLogin 是否登录
type IsLogin struct{}

func (s *IsLogin) Name() string    { return Name }
func (s *IsLogin) Cacheable() bool { return true }
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
