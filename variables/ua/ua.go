package ua

import (
	"context"
	"errors"

	"github.com/liyanbing/filter/cache"
	filterContext "github.com/liyanbing/filter/context"
	"github.com/liyanbing/filter/variables"
)

const Name = "ua"

func init() {
	variables.Register(variables.NewSimpleVariable(&UserAgent{}))
}

// UserAgent 用户代理信息
type UserAgent struct{}

func (s *UserAgent) Name() string    { return Name }
func (s *UserAgent) Cacheable() bool { return true }
func (s *UserAgent) Value(ctx context.Context, _ interface{}, _ *cache.Cache) (interface{}, error) {
	ua, ok := filterContext.FromUA(ctx)
	if !ok {
		return nil, errors.New("ua not found in context")
	}
	return ua, nil
}
