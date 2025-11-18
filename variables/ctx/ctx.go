package ctx

import (
	"context"
	"strings"

	"github.com/liyanbing/filter/cache"
)

const ctxName = "ctx."

func newCtxBuilder() *ctxBuilder {
	return &ctxBuilder{}
}

type ctxBuilder struct{}

func (*ctxBuilder) Name() string {
	return ctxName
}

func (*ctxBuilder) Build(name string) Variable {
	key := strings.TrimPrefix(name, ctxName)
	if key == "" {
		return nil
	}

	return &Ctx{
		name: name,
		key:  key,
	}
}

// Ctx 从上下文的自定义参数中取值
type Ctx struct {
	UnCacheableVariable
	name string
	key  string
}

func (s *Ctx) Name() string { return s.name }
func (s *Ctx) Value(ctx context.Context, _ interface{}, _ *cache.Cache) (interface{}, error) {
	value := ctx.Value(s.key)
	return value, nil
}
