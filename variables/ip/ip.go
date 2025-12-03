package ip

import (
	"context"
	"errors"

	"github.com/airunny/filter/cache"
	filterContext "github.com/airunny/filter/context"
	"github.com/airunny/filter/variables"
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
