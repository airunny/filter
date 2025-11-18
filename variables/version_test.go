package variables

import (
	"context"
	"testing"

	filterContext "github.com/liyanbing/filter/context"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	oldCtx := context.Background()
	ctx := filterContext.WithVersion(oldCtx, _Version)
	ctx2 := filterContext.WithVersion(ctx, "v1.0.0")
	cases := []struct {
		Name  string
		Ctx   context.Context
		Value interface{}
		Error bool
	}{
		{
			Name:  versionName,
			Ctx:   ctx,
			Value: _Version,
		},
		{
			Name:  versionName,
			Ctx:   ctx2,
			Value: "v1.0.0",
		},
		{
			Name:  versionName,
			Ctx:   oldCtx,
			Value: _Version,
			Error: true,
		},
	}

	vv := versionVariable()
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
