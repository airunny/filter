package uid

import (
	"context"
	"errors"

	"github.com/liyanbing/filter/cache"
	filterContext "github.com/liyanbing/filter/context"
	"github.com/liyanbing/filter/variable"
)

const Name = "uid"

func init() {
	variable.Register(variable.NewSimpleVariable(&UID{}))
}

// UID 用户ID
type UID struct{}

func (s *UID) Name() string    { return Name }
func (s *UID) Cacheable() bool { return true }
func (s *UID) Value(ctx context.Context, _ interface{}, _ *cache.Cache) (interface{}, error) {
	uid, ok := filterContext.FromUserId(ctx)
	if !ok {
		return nil, errors.New("uid not found in context")
	}
	return uid, nil
}
