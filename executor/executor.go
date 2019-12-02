package executor

import "context"

type Executor interface {
	Execute(ctx context.Context, data interface{})
}

func BuildExecutor(items []interface{}) (Executor, error) {
	return nil, nil
}
