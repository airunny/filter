package user_tag

import (
	"context"
	"errors"

	"github.com/airunny/filter/cache"
	filterContext "github.com/airunny/filter/context"
	"github.com/airunny/filter/variables"
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
