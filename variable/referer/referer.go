package referer

import (
	"context"
	"errors"

	"github.com/liyanbing/filter/cache"
	filterContext "github.com/liyanbing/filter/context"
	"github.com/liyanbing/filter/variable"
)

const Name = "referer"

func init() {
	variable.Register(variable.NewSimpleVariable(&Referer{}))
}

// Referer referer
type Referer struct{}

func (s *Referer) Name() string    { return Name }
func (s *Referer) Cacheable() bool { return true }
func (s *Referer) Value(ctx context.Context, _ interface{}, _ *cache.Cache) (interface{}, error) {
	value, ok := filterContext.FromReferer(ctx)
	if !ok {
		return nil, errors.New("referer not found in context")
	}
	return value, nil
}
