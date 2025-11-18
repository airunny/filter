package variables

import (
	"context"
	"strings"

	"github.com/liyanbing/calc/compute"
	"github.com/liyanbing/calc/variables"
	"github.com/liyanbing/filter/cache"
	filterType "github.com/liyanbing/filter/filter_type"
)

const calcName = "calc."

func newCalcBuilder() *calcBuilder {
	return &calcBuilder{}
}

type calcBuilder struct{}

func (*calcBuilder) Name() string {
	return calcName
}

func (*calcBuilder) Build(name string) Variable {
	expr := strings.TrimPrefix(name, calcName)
	if expr == "" {
		return nil
	}

	return &Calculator{
		name: name,
		expr: expr,
	}
}

// Calculator 计算器
type Calculator struct {
	UnCacheableVariable
	name string
	expr string
}

func (s *Calculator) Name() string { return s.name }

func (s *Calculator) Value(ctx context.Context, data interface{}, cache *cache.Cache) (interface{}, error) {
	return compute.Evaluate(s.expr, variables.ValueSourceFunc(func(name string) float64 {
		if getter, ok := data.(CalcFactorGetter); ok {
			v, err := getter.CalcFactorGet(ctx, name)
			if err == nil {
				return v
			}
		}

		variable, ok := Get(name)
		if !ok {
			return 0
		}

		value, err := variable.Value(ctx, data, cache)
		if err != nil {
			return 0
		}
		return filterType.GetFloat(value)
	}))
}
