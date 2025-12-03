package time

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/airunny/filter/cache"
	_ "github.com/airunny/filter/location"
	"github.com/airunny/filter/variables"
	"github.com/stretchr/testify/assert"
)

func TestIP(t *testing.T) {
	var (
		ctx = context.Background()
		cc  = cache.NewCache()
		now = time.Now()
	)

	cases := []struct {
		name string
		want func() interface{}
	}{
		{
			name: TimestampName,
			want: func() interface{} {
				return now.Unix()
			},
		},
		{
			name: TsSimpleName,
			want: func() interface{} {
				ret, _ := strconv.ParseUint(now.Format("20060102150405"), 10, 64)
				return ret
			},
		},
		{
			name: SecondName,
			want: func() interface{} {
				return now.Second()
			},
		},
		{
			name: MinuteName,
			want: func() interface{} {
				return now.Minute()
			},
		},
		{
			name: HourName,
			want: func() interface{} {
				return now.Hour()
			},
		},
		{
			name: DayName,
			want: func() interface{} {
				return now.Day()
			},
		},
		{
			name: MonthName,
			want: func() interface{} {
				return int(now.Month())
			},
		},
		{
			name: YearName,
			want: func() interface{} {
				return now.Year()
			},
		},
		{
			name: WdayName,
			want: func() interface{} {
				return int(now.Weekday())
			},
		},
		{
			name: DateName,
			want: func() interface{} {
				return now.Format("2006-01-02")
			},
		},
		{
			name: TimeName,
			want: func() interface{} {
				return now.Format("2006-01-02 15:04:05")
			},
		},
	}

	for index, tt := range cases {
		variable, ok := variables.Get(tt.name)
		assert.True(t, ok)
		assert.NotNil(t, variable)
		assert.Equal(t, tt.name, variable.Name())

		ret, err := variable.Value(ctx, nil, cc)
		assert.Nil(t, err, index)
		assert.Equal(t, tt.want(), ret, index)
	}
}
