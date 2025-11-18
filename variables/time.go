package variables

import (
	"context"
	"strconv"
	"time"

	"github.com/liyanbing/filter/cache"
)

const (
	timestampName = "timestamp"
	tsSimpleName  = "ts_simple"
	secondName    = "second"
	minuteName    = "minute"
	hourName      = "hour"
	dayName       = "day"
	monthName     = "month"
	yearName      = "year"
	wdayName      = "wday"
	dateName      = "date"
	timeName      = "time"
)

func timeStampVariable() *Time {
	return &Time{
		name: timestampName,
	}
}

func tsSimpleVariable() *Time {
	return &Time{
		name: tsSimpleName,
	}
}

func secondVariable() *Time {
	return &Time{
		name: secondName,
	}
}

func minuteVariable() *Time {
	return &Time{
		name: minuteName,
	}
}

func hourVariable() *Time {
	return &Time{
		name: hourName,
	}
}

func dayVariable() *Time {
	return &Time{
		name: dayName,
	}
}

func monthVariable() *Time {
	return &Time{
		name: monthName,
	}
}

func yearVariable() *Time {
	return &Time{
		name: yearName,
	}
}

func wdayVariable() *Time {
	return &Time{
		name: wdayName,
	}
}

func dateVariable() *Time {
	return &Time{
		name: dateName,
	}
}

func timeVariable() *Time {
	return &Time{
		name: timeName,
	}
}

// Time 当前时间的各种形式
type Time struct {
	CacheableVariable
	name string
}

func (s *Time) Name() string { return s.name }
func (s *Time) Value(_ context.Context, _ interface{}, _ *cache.Cache) (interface{}, error) {
	var (
		now   = time.Now()
		value interface{}
	)

	switch s.name {
	case timestampName:
		value = now.Unix()
	case tsSimpleName:
		ret, _ := strconv.ParseUint(now.Format("20060102150405"), 10, 64)
		value = ret
	case secondName:
		value = now.Second()
	case minuteName:
		value = now.Minute()
	case hourName:
		value = now.Hour()
	case dayName:
		value = now.Day()
	case monthName:
		value = int(now.Month())
	case yearName:
		value = now.Year()
	case wdayName:
		value = int(now.Weekday())
	case dateName:
		value = now.Format("2006-01-02")
	default:
		value = now.Format("2006-01-02 15:04:05") // time
	}
	return value, nil
}
