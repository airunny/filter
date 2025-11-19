package context

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithCommonValue(t *testing.T) {
	oldCtx := context.Background()
	ctx := WithUserID(oldCtx, "user_id")
	ctx = WithDevice(ctx, "device_id")
	ctx = WithIP(ctx, "127.0.0.1")
	ctx = WithVersion(ctx, "version")
	ctx = WithPlatform(ctx, "ios")
	ctx = WithChannel(ctx, "channel")
	ctx = WithUA(ctx, "ua")
	ctx = WithReferer(ctx, "referer")
	ctx = WithUserTag(ctx, []string{"1", "2"})

	cases := []struct {
		Expected interface{}
		Got      func(context.Context) (interface{}, bool)
		Exists   bool
	}{
		{
			Expected: "user_id",
			Got: func(ctx context.Context) (interface{}, bool) {
				return FromUserId(ctx)
			},
			Exists: true,
		},
		{
			Expected: nil,
			Got: func(ctx context.Context) (interface{}, bool) {
				return FromUserId(oldCtx)
			},
		},
		{
			Expected: "device_id",
			Got: func(ctx context.Context) (interface{}, bool) {
				return FromDevice(ctx)
			},
			Exists: true,
		},
		{
			Expected: nil,
			Got: func(ctx context.Context) (interface{}, bool) {
				return FromDevice(oldCtx)
			},
		},
		{
			Expected: "127.0.0.1",
			Got: func(ctx context.Context) (interface{}, bool) {
				return FromIP(ctx)
			},
			Exists: true,
		},
		{
			Expected: nil,
			Got: func(ctx context.Context) (interface{}, bool) {
				return FromIP(oldCtx)
			},
		},
		{
			Expected: "version",
			Got: func(ctx context.Context) (interface{}, bool) {
				return FromVersion(ctx)
			},
			Exists: true,
		},
		{
			Expected: nil,
			Got: func(ctx context.Context) (interface{}, bool) {
				return FromVersion(oldCtx)
			},
		},
		{
			Expected: "ios",
			Got: func(ctx context.Context) (interface{}, bool) {
				return FromPlatform(ctx)
			},
			Exists: true,
		},
		{
			Expected: nil,
			Got: func(ctx context.Context) (interface{}, bool) {
				return FromPlatform(oldCtx)
			},
		},
		{
			Expected: "ua",
			Got: func(ctx context.Context) (interface{}, bool) {
				return FromUA(ctx)
			},
			Exists: true,
		},
		{
			Expected: nil,
			Got: func(ctx context.Context) (interface{}, bool) {
				return FromUA(oldCtx)
			},
		},
		{
			Expected: "referer",
			Got: func(ctx context.Context) (interface{}, bool) {
				return FromReferer(ctx)
			},
			Exists: true,
		},
		{
			Expected: nil,
			Got: func(ctx context.Context) (interface{}, bool) {
				return FromReferer(oldCtx)
			},
		},
		{
			Expected: []string{"1", "2"},
			Got: func(ctx context.Context) (interface{}, bool) {
				return FromUserTag(ctx)
			},
			Exists: true,
		},
		{
			Expected: nil,
			Got: func(ctx context.Context) (interface{}, bool) {
				return FromUserTag(oldCtx)
			},
		},
	}

	for _, v := range cases {
		ret, ok := v.Got(ctx)
		assert.Equal(t, v.Exists, ok)
		assert.Equal(t, v.Expected, ret)
	}
}
