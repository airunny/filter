package success

import (
	"context"
	"testing"

	"github.com/liyanbing/filter/cache"
	_ "github.com/liyanbing/filter/location"
	"github.com/liyanbing/filter/variables"
	"github.com/stretchr/testify/assert"
)

func TestReferer(t *testing.T) {
	ctx := context.Background()
	cc := cache.NewCache()
	cases := []struct {
		want interface{}
	}{
		{
			want: Value,
		},
	}

	variable, ok := variables.Get(Name)
	assert.True(t, ok)
	assert.NotNil(t, variable)
	assert.Equal(t, Name, variable.Name())

	for index, tt := range cases {
		ret, err := variable.Value(ctx, nil, cc)
		assert.Nil(t, err)
		assert.Equal(t, tt.want, ret, index)
	}
}
