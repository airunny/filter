package user_tag

import (
	"context"
	"errors"

	"github.com/liyanbing/filter/cache"
	filterContext "github.com/liyanbing/filter/context"
	"github.com/liyanbing/filter/variables"
)

const Name = "user_tag"

func init() {
	variables.Register(variables.NewSimpleVariable(&UserTag{}))
}

type UserTag struct{}

func (s *UserTag) Name() string    { return Name }
func (s *UserTag) Cacheable() bool { return true }
func (s *UserTag) Value(ctx context.Context, _ interface{}, _ *cache.Cache) (interface{}, error) {
	value, ok := filterContext.FromUserTag(ctx)
	if !ok {
		return nil, errors.New("user_tag not found in context")
	}
	return value, nil
}
