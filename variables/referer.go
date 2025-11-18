package variables

import (
	"context"
	"errors"

	"github.com/liyanbing/filter/cache"
	filterContext "github.com/liyanbing/filter/context"
)

const refererName = "referer"

func refererVariable() *Referer {
	return &Referer{}
}

// Referer referer
type Referer struct{ CacheableVariable }

func (s *Referer) Name() string { return refererName }
func (s *Referer) Value(ctx context.Context, _ interface{}, _ *cache.Cache) (interface{}, error) {
	value, ok := filterContext.FromReferer(ctx)
	if !ok {
		return nil, errors.New("referer not found in context")
	}
	return value, nil
}
