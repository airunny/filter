package time

import (
	"context"
	"strconv"
	"time"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/variable"
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
	variable.Register(variable.NewSimpleVariable(&Time{
		name: TimestampName,
	}))
	variable.Register(variable.NewSimpleVariable(&Time{
		name: TsSimpleName,
	}))
	variable.Register(variable.NewSimpleVariable(&Time{
		name: SecondName,
	}))
	variable.Register(variable.NewSimpleVariable(&Time{
		name: MinuteName,
	}))
	variable.Register(variable.NewSimpleVariable(&Time{
		name: HourName,
	}))
	variable.Register(variable.NewSimpleVariable(&Time{
		name: DayName,
	}))
	variable.Register(variable.NewSimpleVariable(&Time{
		name: MonthName,
	}))
	variable.Register(variable.NewSimpleVariable(&Time{
		name: YearName,
	}))
	variable.Register(variable.NewSimpleVariable(&Time{
		name: WdayName,
	}))
	variable.Register(variable.NewSimpleVariable(&Time{
		name: DateName,
	}))
	variable.Register(variable.NewSimpleVariable(&Time{
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
		value = now.Format("2006-01-02 15:04:05") // time
	}
	return value, nil
}
