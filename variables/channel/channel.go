package channel

import (
	"context"
	"errors"

	"github.com/airunny/filter/cache"
	filterContext "github.com/airunny/filter/context"
	"github.com/airunny/filter/variables"
)

const Name = "channel"

func init() {
	variables.Register(variables.NewSimpleVariable(&Channel{}))
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
