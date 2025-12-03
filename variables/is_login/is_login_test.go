package variables

import (
	"context"
	"testing"

	"github.com/airunny/filter/cache"
	filterContext "github.com/airunny/filter/context"
	_ "github.com/airunny/filter/location"
	"github.com/airunny/filter/variables"
	"github.com/stretchr/testify/assert"
)

func TestIsLogin(t *testing.T) {
	ctx := context.Background()
	cc := cache.NewCache()
	cases := []struct {
		ctx  context.Context
		want interface{}
	}{
		{
			ctx:  ctx,
			want: false,
		},
		{
			ctx:  filterContext.WithUserID(ctx, "golang"),
			want: true,
		},
		{
			ctx:  filterContext.WithUserID(ctx, "java"),
			want: true,
		},
		{
			ctx:  filterContext.WithUserID(ctx, "python"),
			want: true,
		},
		{
			ctx:  filterContext.WithUserID(ctx, ""),
			want: false,
		},
		{
			ctx:  filterContext.WithUserID(ctx, "   "),
			want: false,
		},
		{
			ctx:  filterContext.WithUserID(ctx, nil),
			want: false,
		},
		{
			ctx:  filterContext.WithUserID(ctx, 1),
			want: true,
		},
		{
			ctx:  filterContext.WithUserID(ctx, true),
			want: true,
		},
		{
			ctx:  filterContext.WithUserID(ctx, false),
			want: false,
		},
	}

	variable, ok := variables.Get(Name)
	assert.True(t, ok)
	assert.NotNil(t, variable)
	assert.Equal(t, Name, variable.Name())
	for index, tt := range cases {
		ret, err := variable.Value(tt.ctx, nil, cc)
		assert.Nil(t, err)
		assert.Equal(t, tt.want, ret, index)
	}
}
