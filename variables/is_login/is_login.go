package variables

import (
	"context"
	"strings"

	"github.com/airunny/filter/types"

	"github.com/airunny/filter/cache"
	filterContext "github.com/airunny/filter/context"
	"github.com/airunny/filter/variables"
)

const Name = "is_login"

func init() {
	variables.Register(variables.NewSimpleVariable(&IsLogin{}))
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

	if strings.TrimSpace(types.GetString(uid)) != "" {
		return true, nil
	} else {
		return false, nil
	}
}
