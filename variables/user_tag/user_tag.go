package user_tag

import (
	"context"
	"errors"

	"github.com/liyanbing/filter/cache"
	filterContext "github.com/liyanbing/filter/context"
)

const userTabName = "user_tag"

func userTagVariable() *UserTag {
	return &UserTag{}
}

type UserTag struct{ CacheableVariable }

func (s *UserTag) Name() string { return userTabName }
func (s *UserTag) Value(ctx context.Context, _ interface{}, _ *cache.Cache) (interface{}, error) {
	value, ok := filterContext.FromUserTag(ctx)
	if !ok {
		return nil, errors.New("user_tag not found in context")
	}
	return value, nil
}
