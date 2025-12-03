package ctx

import (
	"context"
	"testing"

	"github.com/airunny/filter/cache"
	_ "github.com/airunny/filter/location"
	"github.com/airunny/filter/variables"
	"github.com/stretchr/testify/assert"
)

func TestCtx(t *testing.T) {
	ctx := context.Background()
	cc := cache.NewCache()
	cases := []struct {
		ctx  context.Context
		key  string
		want interface{}
	}{
		// success
		{
			ctx:  ctx,
			key:  "ctx.key",
			want: nil,
		},
		{
			ctx:  context.WithValue(ctx, "key", "tiktok"),
			key:  "ctx.key",
			want: "tiktok",
		},
		{
			ctx:  context.WithValue(ctx, "1", "google"),
			key:  "ctx.1",
			want: "google",
		},
		{
			ctx:  context.WithValue(ctx, "value", 1),
			key:  "ctx.value",
			want: 1,
		},
		{
			ctx:  context.WithValue(ctx, "aa", map[string]interface{}{}),
			key:  "ctx.aa",
			want: map[string]interface{}{},
		},
	}

	for index, tt := range cases {
		variable, ok := variables.Get(tt.key)
		assert.True(t, ok)
		assert.NotNil(t, variable)
		assert.Equal(t, tt.key, variable.Name(), index)

		ret, err := variable.Value(tt.ctx, nil, cc)
		assert.Nil(t, err)
		assert.Equal(t, tt.want, ret, index)
	}
}
