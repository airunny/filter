package variables

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/airunny/filter/cache"
	filterContext "github.com/airunny/filter/context"
	_ "github.com/airunny/filter/location"
	"github.com/airunny/filter/variables"
	"github.com/stretchr/testify/assert"
)

func TestPlatform(t *testing.T) {
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
			err: errors.New("platform not found in context"),
		},
		// success
		{
			ctx:  filterContext.WithPlatform(ctx, "ios"),
			want: "ios",
		},
		{
			ctx:  filterContext.WithPlatform(ctx, "web"),
			want: "web",
		},
		{
			ctx:  filterContext.WithPlatform(ctx, "android"),
			want: "android",
		},
		{
			ctx:  filterContext.WithPlatform(ctx, ""),
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
