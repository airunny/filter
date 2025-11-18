package variables

import "context"

// --------------data getter and setter---------------

type CalcFactorGetter interface {
	CalcFactorGet(ctx context.Context, name string) (float64, error)
}

type FrequencyGetter interface {
	FrequencyGet(ctx context.Context, name string) interface{}
}
