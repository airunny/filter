package filter_context

import (
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithCommonValue(t *testing.T) {
	commonValue := CommonValue{
		UserID:    "110",
		Referer:   "http://www.baidu.com",
		Channel:   "channel",
		UserAgent: "Mozilla/5.0",
		IP:        "127.0.0.1",
		GetForm:   url.Values{"1": []string{"1"}, "2": []string{"2"}},
		Platform:  "ios",
		Device:    "device",
		Version:   "0.0.1",
		UserTags:  []string{"tag1"},
	}

	ctx := context.Background()
	ctx = WithCommonValue(ctx, commonValue)

	cases := []struct {
		Expected interface{}
		Got      func(context.Context) (interface{}, bool)
	}{
		{
			Expected: commonValue.UserID,
			Got: func(ctx context.Context) (interface{}, bool) {
				return UserID(ctx)
			},
		},
		{
			Expected: commonValue.Referer,
			Got: func(ctx context.Context) (interface{}, bool) {
				return Referer(ctx)
			},
		},
		{
			Expected: commonValue.Channel,
			Got: func(ctx context.Context) (interface{}, bool) {
				return Channel(ctx)
			},
		},
		{
			Expected: commonValue.UserAgent,
			Got: func(ctx context.Context) (interface{}, bool) {
				return UserAgent(ctx)
			},
		},
		{
			Expected: commonValue.IP,
			Got: func(ctx context.Context) (interface{}, bool) {
				return IP(ctx)
			},
		},
		{
			Expected: commonValue.GetForm,
			Got: func(ctx context.Context) (interface{}, bool) {
				return Form(ctx)
			},
		},
		{
			Expected: commonValue.Platform,
			Got: func(ctx context.Context) (interface{}, bool) {
				return Platform(ctx)
			},
		},
		{
			Expected: commonValue.Device,
			Got: func(ctx context.Context) (interface{}, bool) {
				return Device(ctx)
			},
		},
		{
			Expected: commonValue.Version,
			Got: func(ctx context.Context) (interface{}, bool) {
				return Version(ctx)
			},
		},
		{
			Expected: commonValue.UserTags,
			Got: func(ctx context.Context) (interface{}, bool) {
				return UserTags(ctx)
			},
		},
	}

	for _, v := range cases {
		ret, ok := v.Got(ctx)
		assert.Equal(t, true, ok)
		assert.Equal(t, v.Expected, ret)
	}
}

func TestWithCustom(t *testing.T) {
	custom := map[string]interface{}{
		"1": 1,
		"2": 2,
		"test": struct {
		}{},
		"4": true,
	}
	ctx := context.Background()
	ctx = WithCustom(ctx, custom)

	ret, ok := FromCustom(ctx)
	assert.Equal(t, true, ok)
	assert.Equal(t, custom, ret)
}
