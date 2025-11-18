package success

import (
	"context"

	"github.com/liyanbing/filter/cache"
)

const (
	successName  = "success"
	successValue = 1
)

func successVariable() *Success {
	return &Success{}
}

// Success 永远返回1
type Success struct{ CacheableVariable }

func (s *Success) Name() string { return successName }

func (s *Success) Value(_ context.Context, _ interface{}, _ *cache.Cache) (interface{}, error) {
	return successValue, nil
}
