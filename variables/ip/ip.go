package ip

import (
	"context"
	"errors"

	"github.com/liyanbing/filter/cache"
	filterContext "github.com/liyanbing/filter/context"
	"github.com/liyanbing/filter/variables"
)

const Name = "ip"

func init() {
	variables.Register(variables.NewSimpleVariable(&IP{}))
}

// IP 从上下文中获取IP地址
type IP struct{}

func (s *IP) Name() string    { return Name }
func (s *IP) Cacheable() bool { return true }
func (s *IP) Value(ctx context.Context, _ interface{}, _ *cache.Cache) (interface{}, error) {
	ip, ok := filterContext.FromIP(ctx)
	if !ok {
		return nil, errors.New("ip not found in context")
	}
	return ip, nil
}
