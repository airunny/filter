package ip

import (
	"context"
	"errors"

	"github.com/liyanbing/filter/cache"
	filterContext "github.com/liyanbing/filter/context"
)

const IPName = "ip"

func ipVariable() *IP {
	return &IP{}
}

// IP 从上下文中获取IP地址
type IP struct{ CacheableVariable }

func (s *IP) Name() string { return IPName }

func (s *IP) Value(ctx context.Context, _ interface{}, _ *cache.Cache) (interface{}, error) {
	ip, ok := filterContext.FromIP(ctx)
	if !ok {
		return nil, errors.New("ip not found in context")
	}
	return ip, nil
}
