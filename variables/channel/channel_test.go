package channel

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
			err: errors.New("channel not found in context"),
		},
		// success
		{
			ctx:  filterContext.WithChannel(ctx, "weichat"),
			want: "weichat",
		},
		{
			ctx:  filterContext.WithChannel(ctx, "tiktok"),
			want: "tiktok",
		},
		{
			ctx:  filterContext.WithChannel(ctx, "google"),
			want: "google",
		},
		{
			ctx:  filterContext.WithChannel(ctx, ""),
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
