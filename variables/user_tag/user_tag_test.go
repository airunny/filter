package user_tag

import (
	"context"
	"testing"

	filterContext "github.com/liyanbing/filter/context"
	"github.com/liyanbing/filter/variables"
	"github.com/stretchr/testify/assert"
)

func Test_UserTag(t *testing.T) {
	oldCtx := context.Background()
	ctx := filterContext.WithUserTag(oldCtx, []string{"1", "2"})
	ctx2 := filterContext.WithUserTag(oldCtx, []string{"3", "4"})
	cases := []struct {
		Name  string
		Ctx   context.Context
		Value interface{}
		Error bool
	}{
		{
			Name:  Name,
			Ctx:   ctx,
			Value: []string{"1", "2"},
		},
		{
			Name:  Name,
			Ctx:   ctx2,
			Value: []string{"3", "4"},
		},
		{
			Name:  Name,
			Ctx:   oldCtx,
			Value: nil,
			Error: true,
		},
	}

	vv, exists := variables.Get(Name)
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
