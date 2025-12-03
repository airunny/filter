package uid

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

func TestUID(t *testing.T) {
	ctx := context.Background()
	cc := cache.NewCache()
	cases := []struct {
		ctx  context.Context
		want interface{}
		err  error
	}{
		// err
		{
			ctx: ctx,
			err: errors.New("uid not found in context"),
		},
		// success
		{
			ctx:  filterContext.WithUserID(ctx, "golang"),
			want: "golang",
		},
		{
			ctx:  filterContext.WithUserID(ctx, "java"),
			want: "java",
		},
		{
			ctx:  filterContext.WithUserID(ctx, "python"),
			want: "python",
		},
		{
			ctx:  filterContext.WithUserID(ctx, ""),
			want: "",
		},
		{
			ctx:  filterContext.WithUserID(ctx, true),
			want: true,
		},
		{
			ctx:  filterContext.WithUserID(ctx, false),
			want: false,
		},
		{
			ctx:  filterContext.WithUserID(ctx, 1),
			want: 1,
		},
		{
			ctx:  filterContext.WithUserID(ctx, 10.10),
			want: 10.10,
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
