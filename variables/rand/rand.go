package variables

import (
	"context"
	"math/rand"

	"github.com/airunny/filter/cache"
	"github.com/airunny/filter/variables"
)

const Name = "rand"

func init() {
	variables.Register(variables.NewSimpleVariable(&Rand{}))
}

// Rand 随机返回1-100之间的值[1,100]
type Rand struct{}

func (s *Rand) Name() string    { return Name }
func (s *Rand) Cacheable() bool { return false }
func (s *Rand) Value(_ context.Context, _ interface{}, _ *cache.Cache) (interface{}, error) {
	return rand.Intn(100) + 1, nil
}
