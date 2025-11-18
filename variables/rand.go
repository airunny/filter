package variables

import (
	"context"
	"math/rand"

	"github.com/liyanbing/filter/cache"
)

const randName = "rand"

func randVariable() *Rand {
	return &Rand{}
}

// Rand 随机返回1-100之间的值[1,100]
type Rand struct{ UnCacheableVariable }

func (s *Rand) Name() string { return randName }

func (s *Rand) Value(_ context.Context, _ interface{}, _ *cache.Cache) (interface{}, error) {
	return rand.Intn(100) + 1, nil
}
