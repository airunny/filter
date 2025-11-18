package variables

import "context"

// --------------data getter and setter---------------

type CalcFactorGetter interface {
	CalcValue(ctx context.Context, name string) (float64, error)
}

type FrequencyGetter interface {
	FrequencyValue(ctx context.Context, name string) interface{}
}
