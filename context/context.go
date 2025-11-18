package context

import (
	"context"
)

type (
	userIdKey   struct{}
	deviceKey   struct{}
	ipKey       struct{}
	versionKey  struct{}
	platformKey struct{}
	channelKey  struct{}
	uaKey       struct{}
	refererKey  struct{}
	userTagKey  struct{}
)

func WithUserID(ctx context.Context, userId interface{}) context.Context {
	return context.WithValue(ctx, userIdKey{}, userId)
}
func FromUserId(ctx context.Context) (interface{}, bool) {
	value, ok := ctx.Value(userIdKey{}).(interface{})
	return value, ok
}

func WithDevice(ctx context.Context, device interface{}) context.Context {
	return context.WithValue(ctx, deviceKey{}, device)
}
func FromDevice(ctx context.Context) (interface{}, bool) {
	value, ok := ctx.Value(deviceKey{}).(interface{})
	return value, ok
}

func WithIP(ctx context.Context, ip interface{}) context.Context {
	return context.WithValue(ctx, ipKey{}, ip)
}
func FromIP(ctx context.Context) (interface{}, bool) {
	value, ok := ctx.Value(ipKey{}).(interface{})
	return value, ok
}

func WithVersion(ctx context.Context, version interface{}) context.Context {
	return context.WithValue(ctx, versionKey{}, version)
}
func FromVersion(ctx context.Context) (interface{}, bool) {
	value, ok := ctx.Value(versionKey{}).(interface{})
	return value, ok
}

func WithPlatform(ctx context.Context, platform interface{}) context.Context {
	return context.WithValue(ctx, platformKey{}, platform)
}
func FromPlatform(ctx context.Context) (interface{}, bool) {
	value, ok := ctx.Value(platformKey{}).(interface{})
	return value, ok
}

func WithChannel(ctx context.Context, channel interface{}) context.Context {
	return context.WithValue(ctx, channelKey{}, channel)
}
func FromChannel(ctx context.Context) (interface{}, bool) {
	value, ok := ctx.Value(channelKey{}).(interface{})
	return value, ok
}

func WithUA(ctx context.Context, ua interface{}) context.Context {
	return context.WithValue(ctx, uaKey{}, ua)
}
func FromUA(ctx context.Context) (interface{}, bool) {
	value, ok := ctx.Value(uaKey{}).(interface{})
	return value, ok
}

func WithReferer(ctx context.Context, referer interface{}) context.Context {
	return context.WithValue(ctx, refererKey{}, referer)
}
func FromReferer(ctx context.Context) (interface{}, bool) {
	value, ok := ctx.Value(refererKey{}).(interface{})
	return value, ok
}

func WithUserTag(ctx context.Context, userTag interface{}) context.Context {
	return context.WithValue(ctx, userTagKey{}, userTag)
}
func FromUserTag(ctx context.Context) (interface{}, bool) {
	value, ok := ctx.Value(userTagKey{}).(interface{})
	return value, ok
}
