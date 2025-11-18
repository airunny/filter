package variables

import (
	"context"
	"errors"

	"github.com/liyanbing/filter/cache"
	filterContext "github.com/liyanbing/filter/context"
)

const uidName = "uid"

func uidVariable() *UID {
	return &UID{}
}

// UID 用户ID
type UID struct{ CacheableVariable }

func (s *UID) Name() string { return uidName }
func (s *UID) Value(ctx context.Context, _ interface{}, _ *cache.Cache) (interface{}, error) {
	uid, ok := filterContext.FromUserId(ctx)
	if !ok {
		return nil, errors.New("uid not found in context")
	}
	return uid, nil
}
