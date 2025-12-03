package time

import (
	"context"
	"strconv"
	"time"

	"github.com/airunny/filter/cache"
	"github.com/airunny/filter/variables"
)

const (
	TimestampName = "timestamp"
	TsSimpleName  = "ts_simple"
	SecondName    = "second"
	MinuteName    = "minute"
	HourName      = "hour"
	DayName       = "day"
	MonthName     = "month"
	YearName      = "year"
	WdayName      = "wday"
	DateName      = "date"
	TimeName      = "time"
)

func init() {
	variables.Register(variables.NewSimpleVariable(&Time{
		name: TimestampName,
	}))
	variables.Register(variables.NewSimpleVariable(&Time{
		name: TsSimpleName,
	}))
	variables.Register(variables.NewSimpleVariable(&Time{
		name: SecondName,
	}))
	variables.Register(variables.NewSimpleVariable(&Time{
		name: MinuteName,
	}))
	variables.Register(variables.NewSimpleVariable(&Time{
		name: HourName,
	}))
	variables.Register(variables.NewSimpleVariable(&Time{
		name: DayName,
	}))
	variables.Register(variables.NewSimpleVariable(&Time{
		name: MonthName,
	}))
	variables.Register(variables.NewSimpleVariable(&Time{
		name: YearName,
	}))
	variables.Register(variables.NewSimpleVariable(&Time{
		name: WdayName,
	}))
	variables.Register(variables.NewSimpleVariable(&Time{
		name: DateName,
	}))
	variables.Register(variables.NewSimpleVariable(&Time{
		name: TimeName,
	}))

}

// Time 当前时间的各种形式
type Time struct {
	name string
}

func (s *Time) Name() string    { return s.name }
func (s *Time) Cacheable() bool { return false }
func (s *Time) Value(_ context.Context, _ interface{}, _ *cache.Cache) (interface{}, error) {
	var (
		now   = time.Now()
		value interface{}
	)

	switch s.name {
	case TimestampName:
		value = now.Unix()
	case TsSimpleName:
		ret, _ := strconv.ParseUint(now.Format("20060102150405"), 10, 64)
		value = ret
	case SecondName:
		value = now.Second()
	case MinuteName:
		value = now.Minute()
	case HourName:
		value = now.Hour()
	case DayName:
		value = now.Day()
	case MonthName:
		value = int(now.Month())
	case YearName:
		value = now.Year()
	case WdayName:
		value = int(now.Weekday())
	case DateName:
		value = now.Format("2006-01-02")
	default:
		value = now.Format("2006-01-02 15:04:05")
	}
	return value, nil
}
