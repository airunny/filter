package ip

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

func TestIP(t *testing.T) {
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
			err: errors.New("ip not found in context"),
		},
		// success
		{
			ctx:  filterContext.WithIP(ctx, "47.107.69.99"),
			want: "47.107.69.99",
		},
		{
			ctx:  filterContext.WithIP(ctx, "127.0.0.1"),
			want: "127.0.0.1",
		},
		{
			ctx:  filterContext.WithIP(ctx, "192.168.1.1"),
			want: "192.168.1.1",
		},
		{
			ctx:  filterContext.WithIP(ctx, ""),
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
