package success

import (
	"context"

	"github.com/airunny/filter/cache"
	"github.com/airunny/filter/variables"
)

const (
	Name  = "success"
	Value = 1
)

func init() {
	variables.Register(variables.NewSimpleVariable(&Success{}))
}

// Success 永远返回1
type Success struct{}

func (s *Success) Name() string    { return Name }
func (s *Success) Cacheable() bool { return true }
func (s *Success) Value(_ context.Context, _ interface{}, _ *cache.Cache) (interface{}, error) {
	return Value, nil
}
