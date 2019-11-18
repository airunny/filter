package filter

import (
	"context"
	"encoding/json"
	"sync/atomic"
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
type CommonConf struct {
	JsonStr string
}

type base struct {
	val atomic.Value
}

type baseValuePair struct {
	Cfg      *baseCfg
	Reporter Reporter
}

type baseCfg struct {
	M map[string]baseValues `json:"m"`
	V string                `json:"version"`
}

type baseValues struct {
	FilterData json.RawMessage `json:"filter_data"`
	Weight     int64           `json:"weight"`
	Priority   int64           `json:"priority"`
}
