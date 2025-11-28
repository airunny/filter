package device

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/liyanbing/filter/cache"
	filterContext "github.com/liyanbing/filter/context"
	_ "github.com/liyanbing/filter/location"
	"github.com/liyanbing/filter/variables"
	"github.com/stretchr/testify/assert"
)

func TestDevice(t *testing.T) {
	ctx := context.Background()
	cc := cache.NewCache()
	cases := []struct {
		ctx  context.Context
		want string
		err  error
	}{
		// err
		{
			ctx: ctx,
			err: errors.New("device not found in context"),
		},
		// success
		{
			ctx:  filterContext.WithDevice(ctx, "1"),
			want: "1",
		},
		{
			ctx:  filterContext.WithDevice(ctx, "2"),
			want: "2",
		},
		{
			ctx:  filterContext.WithDevice(ctx, "3"),
			want: "3",
		},
		{
			ctx:  filterContext.WithDevice(ctx, ""),
			want: "",
		},
	}

	variable, ok := variables.Get(Name)
	assert.True(t, ok)
	assert.NotNil(t, variable)
	assert.Equal(t, Name, variable.Name())
	for index, tt := range cases {
		ret, err := variable.Value(tt.ctx, nil, cc)
		if err != nil {
			assert.True(t, reflect.DeepEqual(tt.err, err), index)
		} else {
			assert.Equal(t, tt.want, ret, index)
		}
	}
}
