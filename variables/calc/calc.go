package calc

import (
	"context"
	"strings"

	"github.com/airunny/filter/cache"
	"github.com/airunny/filter/variables"
	"github.com/liyanbing/calc/compute"
	calcVariables "github.com/liyanbing/calc/variables"
)

const Name = "calc."

func init() {
	variables.Register(&calcBuilder{})
}

type calcBuilder struct{}

func (*calcBuilder) Name() string {
	return Name
}

func (*calcBuilder) Build(name string) variables.Variable {
	expr := strings.TrimPrefix(name, Name)
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
	name string
	expr string
}

func (s *Calculator) Cacheable() bool { return false }
func (s *Calculator) Name() string    { return s.name }
func (s *Calculator) Value(ctx context.Context, data interface{}, cache *cache.Cache) (interface{}, error) {
	return compute.Evaluate(s.expr, calcVariables.ValueSourceFunc(func(key string) float64 {
		if getter, ok := data.(variables.Calculator); ok {
			v, err := getter.CalcValue(ctx, key)
			if err == nil {
				return v
			}
		}
		return 0
		//vv, ok := variables.Get(key)
		//if !ok {
		//	return 0
		//}
		//
		//value, err := vv.Value(ctx, data, cache)
		//if err != nil {
		//	return 0
		//}
		//return types.GetFloat(value)
	}))
}
