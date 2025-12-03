package calc

import (
	"context"
	"reflect"
	"strconv"
	"testing"

	"github.com/airunny/filter/cache"
	_ "github.com/airunny/filter/location"
	"github.com/airunny/filter/variables"
	_ "github.com/airunny/filter/variables/ctx"
	"github.com/stretchr/testify/assert"
)

type CalcValue struct {
	err error
}

func (s *CalcValue) CalcValue(ctx context.Context, key string) (float64, error) {
	if s.err != nil {
		return 0, s.err
	}

	switch key {
	case "name":
		return float64(10), nil
	case "age":
		return float64(20), nil
	}
	return strconv.ParseFloat(key, 64)
}

func TestCalc(t *testing.T) {
	ctx := context.Background()
	cc := cache.NewCache()
	cases := []struct {
		data    interface{}
		express string
		want    interface{}
		err     error
	}{
		{
			data:    &CalcValue{},
			express: "calc.1 + 2",
			want:    float64(3),
		},
		{
			data:    &CalcValue{},
			express: "calc.2 + 1",
			want:    float64(3),
		},
		{
			data:    &CalcValue{},
			express: "calc.1-2",
			want:    float64(-1),
		},
		{
			data:    &CalcValue{},
			express: "calc.2-1",
			want:    float64(1),
		},
		{
			data:    &CalcValue{},
			express: "calc.1*2",
			want:    float64(2),
		},
		{
			data:    &CalcValue{},
			express: "calc.2*1",
			want:    float64(2),
		},
		{
			data:    &CalcValue{},
			express: "calc.1/2",
			want:    0.5,
		},
		{
			data:    &CalcValue{},
			express: "calc.2/1",
			want:    float64(2),
		},
		{
			data:    &CalcValue{},
			express: "calc.__name*__age",
			want:    float64(200),
		},
	}

	for index, tt := range cases {
		variable, ok := variables.Get(tt.express)
		assert.True(t, ok)
		assert.NotNil(t, variable)
		assert.Equal(t, tt.express, variable.Name())

		ret, err := variable.Value(ctx, tt.data, cc)
		if err != nil {
			assert.True(t, reflect.DeepEqual(tt.err, err), index)
		} else {
			assert.Equal(t, tt.want, ret, index)
		}
	}
}
