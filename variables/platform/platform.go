package variables

import (
	"context"
	"errors"

	"github.com/airunny/filter/cache"
	filterContext "github.com/airunny/filter/context"
	"github.com/airunny/filter/variables"
)

const Name = "platform"

func init() {
	variables.Register(variables.NewSimpleVariable(&Platform{}))
}

// Platform 平台
type Platform struct{}

func (s *Platform) Name() string    { return Name }
func (s *Platform) Cacheable() bool { return true }
func (s *Platform) Value(ctx context.Context, _ interface{}, _ *cache.Cache) (interface{}, error) {
	plt, ok := filterContext.FromPlatform(ctx)
	if !ok {
		return nil, errors.New("platform not found in context")
	}
	return plt, nil
}
