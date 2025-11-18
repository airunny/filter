package variables

import (
	"context"
	"errors"
	"regexp"

	"github.com/liyanbing/filter/cache"
)

var getReg = regexp.MustCompile(`^get.(.+)`)

type Variable interface {
	Name() string
	Cacheable() bool
	Value(ctx context.Context, data interface{}, cache *cache.Cache) (interface{}, error)
}

type Builder interface {
	Name() string
	Build(string) Variable
}

type CacheableVariable struct{}

func (s *CacheableVariable) Cacheable() bool {
	return true
}

type UnCacheableVariable struct{}

func (s *UnCacheableVariable) Cacheable() bool {
	return false
}

func GetVariableValue(ctx context.Context, v Variable, data interface{}, cache *cache.Cache) (interface{}, error) {
	if v == nil {
		return nil, errors.New("empty variable")
	}

	if v.Cacheable() {
		if value, ok := cache.Get(v.Name()); ok {
			return value, nil
		}
	}

	value, err := v.Value(ctx, data, cache)
	if err != nil {
		return nil, err
	}

	if v.Cacheable() {
		cache.Set(v.Name(), value)
	}
	return value, nil
}
