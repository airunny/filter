package version

import (
	"context"
	"testing"

	filterContext "github.com/liyanbing/filter/context"
	"github.com/liyanbing/filter/variable"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	oldCtx := context.Background()
	ctx := filterContext.WithVersion(oldCtx, "version")
	ctx2 := filterContext.WithVersion(ctx, "v1.0.0")
	cases := []struct {
		Name  string
		Ctx   context.Context
		Value interface{}
		Error bool
	}{
		{
			Name:  Name,
			Ctx:   ctx,
			Value: "version",
		},
		{
			Name:  Name,
			Ctx:   ctx2,
			Value: "v1.0.0",
		},
		{
			Name:  Name,
			Ctx:   oldCtx,
			Value: "version",
			Error: true,
		},
	}

	vv, exists := variable.Get(Name)
	assert.True(t, exists)
	assert.NotNil(t, vv)
	for _, tt := range cases {
		assert.Equal(t, vv.Name(), tt.Name)
		value, err := vv.Value(tt.Ctx, nil, nil)
		if tt.Error {
			assert.NotNil(t, err)
		} else {
			assert.Equal(t, tt.Value, value)
		}
	}
}
