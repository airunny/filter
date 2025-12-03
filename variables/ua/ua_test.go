package ua

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

func TestChannel(t *testing.T) {
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
			err: errors.New("ua not found in context"),
		},
		// success
		{
			ctx:  filterContext.WithUA(ctx, "ua1"),
			want: "ua1",
		},
		{
			ctx:  filterContext.WithUA(ctx, "ua2"),
			want: "ua2",
		},
		{
			ctx:  filterContext.WithUA(ctx, "ua3"),
			want: "ua3",
		},
		{
			ctx:  filterContext.WithUA(ctx, ""),
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
