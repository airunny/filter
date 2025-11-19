package success

import (
	"context"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/variable"
)

const (
	successName  = "success"
	successValue = 1
)

func init() {
	variable.Register(variable.NewSimpleVariable(&Success{}))
}

// Success 永远返回1
type Success struct{}

func (s *Success) Name() string    { return successName }
func (s *Success) Cacheable() bool { return true }
func (s *Success) Value(_ context.Context, _ interface{}, _ *cache.Cache) (interface{}, error) {
	return successValue, nil
}
