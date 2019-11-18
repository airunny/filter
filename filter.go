package filter

import "context"

const filterContextVariableKey = "filter.ctx.variable"

func WithContext(ctx context.Context, data map[string]interface{}) context.Context {
	return context.WithValue(ctx, filterContextVariableKey, data)
}

func FromContext(ctx context.Context) (map[string]interface{}, bool) {
	data := ctx.Value(filterContextVariableKey)
	if value, ok := data.(map[string]interface{}); ok {
		return value, true
	}
	return nil, false
}
