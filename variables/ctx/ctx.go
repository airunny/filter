package ctx

import (
	"context"
	"strings"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/variables"
)

const Name = "ctx."

func init() {
	variables.Register(&ctxBuilder{})
}

type ctxBuilder struct{}

func (*ctxBuilder) Name() string {
	return Name
}

func (*ctxBuilder) Build(name string) variables.Variable {
	key := strings.TrimPrefix(name, Name)
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
	name string
	key  string
}

func (s *Ctx) Name() string    { return s.name }
func (s *Ctx) Cacheable() bool { return false }
func (s *Ctx) Value(ctx context.Context, _ interface{}, _ *cache.Cache) (interface{}, error) {
	value := ctx.Value(s.key)
	return value, nil
}
