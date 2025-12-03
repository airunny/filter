package variables

import (
	"context"
	"testing"

	"github.com/airunny/filter/cache"
	_ "github.com/airunny/filter/location"
	"github.com/airunny/filter/variables"
	_ "github.com/airunny/filter/variables/ip"
	"github.com/stretchr/testify/assert"
)

func TestArea(t *testing.T) {
	ctx := context.Background()
	cc := cache.NewCache()

	variable, ok := variables.Get(Name)
	assert.True(t, ok)
	assert.NotNil(t, variable)
	assert.Equal(t, Name, variable.Name())

	for i := 0; i < 1000000; i++ {
		ret, err := variable.Value(ctx, nil, cc)
		assert.Nil(t, err)
		assert.GreaterOrEqual(t, ret, 1)
		assert.LessOrEqual(t, ret, 100)
	}
}
