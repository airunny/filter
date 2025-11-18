package channel

import (
	"context"
	"errors"

	"github.com/liyanbing/filter/cache"
	filterContext "github.com/liyanbing/filter/context"
	"github.com/liyanbing/filter/variables"
)

const channelName = "channel"

func channelVariable() *Channel {
	return &Channel{}
}

// Channel 渠道
type Channel struct{ variables.CacheableVariable }

func (s *Channel) Name() string { return channelName }
func (s *Channel) Value(ctx context.Context, _ interface{}, _ *cache.Cache) (interface{}, error) {
	value, ok := filterContext.FromChannel(ctx)
	if !ok {
		return nil, errors.New("channel not found in context")
	}
	return value, nil
}
