package version

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
			err: errors.New("version not found in context"),
		},
		// success
		{
			ctx:  filterContext.WithVersion(ctx, "v1.0.0"),
			want: "v1.0.0",
		},
		{
			ctx:  filterContext.WithVersion(ctx, "v2.0.0"),
			want: "v2.0.0",
		},
		{
			ctx:  filterContext.WithVersion(ctx, "v3.0.0"),
			want: "v3.0.0",
		},
		{
			ctx:  filterContext.WithVersion(ctx, ""),
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
