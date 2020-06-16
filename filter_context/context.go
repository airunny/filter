package filter_context

import (
	"context"
	"net/url"
)

type contextValueKeyType int

const (
	filterContextGeneralVariableKey contextValueKeyType = iota
	filterContextCustomVariableKey
)

type CommonValue struct {
	UserID    string     // 用户ID
	Referer   string     // referer
	Channel   string     // 渠道
	UserAgent string     // user agent
	IP        string     // ip
	GetForm   url.Values // 请求体
	Platform  string     // 平台
	Device    string     // 设备
	Version   string     // 版本
	UserTags  []string   // 用户标签
}

// ----------- common values
func WithCommonValue(parent context.Context, value CommonValue) context.Context {
	return context.WithValue(parent, filterContextGeneralVariableKey, value)
}

func FromCommonValue(ctx context.Context) (CommonValue, bool) {
	val, ok := ctx.Value(filterContextGeneralVariableKey).(CommonValue)
	return val, ok
}

// ------------ custom values
func WithCustom(ctx context.Context, data map[string]interface{}) context.Context {
	return context.WithValue(ctx, filterContextCustomVariableKey, data)
}

func FromCustom(ctx context.Context) (map[string]interface{}, bool) {
	data := ctx.Value(filterContextCustomVariableKey)
	if value, ok := data.(map[string]interface{}); ok {
		return value, true
	}

	return nil, false
}

// ----------
func UserID(ctx context.Context) (string, bool) {
	val, ok := FromCommonValue(ctx)
	if !ok {
		return "", false
	}

	return val.UserID, val.UserID != ""
}

func Referer(ctx context.Context) (string, bool) {
	val, ok := FromCommonValue(ctx)
	if !ok {
		return "", false
	}

	return val.Referer, val.Referer != ""
}

func Channel(ctx context.Context) (string, bool) {
	val, ok := FromCommonValue(ctx)
	if !ok {
		return "", false
	}

	return val.Channel, val.Channel != ""
}

func UserAgent(ctx context.Context) (string, bool) {
	val, ok := FromCommonValue(ctx)
	if !ok {
		return "", false
	}

	return val.UserAgent, val.UserAgent != ""
}

func IP(ctx context.Context) (string, bool) {
	val, ok := FromCommonValue(ctx)
	if !ok {
		return "", false
	}

	return val.IP, val.IP != ""
}

func Form(ctx context.Context) (url.Values, bool) {
	val, ok := FromCommonValue(ctx)
	if !ok {
		return nil, false
	}

	return val.GetForm, val.GetForm != nil
}

func Platform(ctx context.Context) (string, bool) {
	val, ok := FromCommonValue(ctx)
	if !ok {
		return "", false
	}

	return val.Platform, val.Platform != ""
}

func Device(ctx context.Context) (string, bool) {
	val, ok := FromCommonValue(ctx)
	if !ok {
		return "", false
	}

	return val.Device, val.Device != ""
}

func Version(ctx context.Context) (string, bool) {
	val, ok := FromCommonValue(ctx)
	if !ok {
		return "", false
	}

	return val.Version, val.Version != ""
}

func UserTags(ctx context.Context) ([]string, bool) {
	val, ok := FromCommonValue(ctx)
	if !ok {
		return nil, false
	}

	return val.UserTags, val.UserTags != nil
}
