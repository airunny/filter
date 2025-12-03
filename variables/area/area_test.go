package area

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/airunny/filter/cache"
	filterContext "github.com/airunny/filter/context"
	_ "github.com/airunny/filter/location"
	"github.com/airunny/filter/variables"
	_ "github.com/airunny/filter/variables/ip"
	"github.com/stretchr/testify/assert"
)

func TestArea(t *testing.T) {
	ctx := context.Background()
	ipCtx := filterContext.WithIP(ctx, "47.107.69.99")
	cc := cache.NewCache()

	cases := []struct {
		ctx  context.Context
		name string
		want string
		err  error
	}{
		// err
		{
			name: CityName,
			ctx:  ctx,
			err:  errors.New("ip not found in context"),
		},
		// success
		{
			ctx:  ipCtx,
			name: CountryName,
			want: "China",
		},
		{
			ctx:  ipCtx,
			name: ProvinceName,
			want: "Guangdong",
		},
		{
			ctx:  ipCtx,
			name: CityName,
			want: "Shenzhen",
		},
	}

	for _, tt := range cases {
		variable, ok := variables.Get(tt.name)
		assert.True(t, ok)
		assert.NotNil(t, variable)
		assert.Equal(t, tt.name, variable.Name())

		ret, err := variables.GetValue(tt.ctx, variable, nil, cc)
		if err != nil {
			assert.True(t, reflect.DeepEqual(tt.err, err))
		} else {
			assert.Equal(t, tt.want, ret)
		}
	}
}
