package channel

import (
	"context"
	"errors"

	"github.com/liyanbing/filter/cache"
	filterContext "github.com/liyanbing/filter/context"
	"github.com/liyanbing/filter/variable"
)

const Name = "channel"

func init() {
	variable.Register(variable.NewSimpleVariable(&Channel{}))
}

// Channel 渠道
type Channel struct{}

func (s *Channel) Name() string    { return Name }
func (s *Channel) Cacheable() bool { return true }
func (s *Channel) Value(ctx context.Context, _ interface{}, _ *cache.Cache) (interface{}, error) {
	value, ok := filterContext.FromChannel(ctx)
	if !ok {
		return nil, errors.New("channel not found in context")
	}
	return value, nil
}
