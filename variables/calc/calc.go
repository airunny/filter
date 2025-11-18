package calc

import (
	"context"
	"strings"

	"github.com/liyanbing/calc/compute"
	calcVariables "github.com/liyanbing/calc/variables"
	"github.com/liyanbing/filter/cache"
	filterType "github.com/liyanbing/filter/filter_type"
	"github.com/liyanbing/filter/variables"
)

const calcName = "calc."

func newCalcBuilder() *calcBuilder {
	return &calcBuilder{}
}

type calcBuilder struct{}

func (*calcBuilder) Name() string {
	return calcName
}

func (*calcBuilder) Build(name string) *Calculator {
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
	variables.UnCacheableVariable
	name string
	expr string
}

func (s *Calculator) Name() string { return s.name }

func (s *Calculator) Value(ctx context.Context, data interface{}, cache *cache.Cache) (interface{}, error) {
	return compute.Evaluate(s.expr, calcVariables.ValueSourceFunc(func(name string) float64 {
		if getter, ok := data.(variables.CalcFactorGetter); ok {
			v, err := getter.CalcValue(ctx, name)
			if err == nil {
				return v
			}
		}

		variable, ok := variables.Get(name)
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
