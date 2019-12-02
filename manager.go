package filter

import (
	"context"
)

type Manger interface {
	Execute(ctx context.Context, data interface{}) (ret interface{}, err error)
	Refresh(ctx context.Context, jsonStr string) error
}

type Reporter interface {
	Report(ctx context.Context, version string, succ int, succId string, data interface{})
}

type ReportFunc func(ctx context.Context, version string, succ int, succId string, data interface{})

func (rf ReportFunc) Report(ctx context.Context, version string, succ int, succId string, data interface{}) {
	rf(ctx, version, succ, succId, data)
}

// -------------
type Config struct {
	Filters map[string]SingleConfig `json:"filters"`
	Version string                  `json:"version"`
}

type SingleConfig struct {
	FilterData []interface{} `json:"filter_data"`
	Weight     int64         `json:"weight"`
	Priority   int64         `json:"priority"`
}
