package user_tag

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

func TestUserTag(t *testing.T) {
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
			err: errors.New("user_tag not found in context"),
		},
		// success
		{
			ctx:  filterContext.WithUserTag(ctx, []string{"1", "2"}),
			want: []string{"1", "2"},
		},
		{
			ctx:  filterContext.WithUserTag(ctx, []string{"3", "4"}),
			want: []string{"3", "4"},
		},
		{
			ctx:  filterContext.WithUserTag(ctx, []string{"1"}),
			want: []string{"1"},
		},
		{
			ctx:  filterContext.WithUserTag(ctx, []string{}),
			want: []string{},
		},
		{
			ctx:  filterContext.WithUserTag(ctx, "user_tag"),
			want: "user_tag",
		},
		{
			ctx:  filterContext.WithUserTag(ctx, ""),
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
