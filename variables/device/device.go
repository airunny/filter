package device

import (
	"context"
	"errors"

	"github.com/liyanbing/filter/cache"
	filterContext "github.com/liyanbing/filter/context"
)

const deviceName = "device"

func deviceVariable() *Device {
	return &Device{}
}

// Device 设备ID
type Device struct{ CacheableVariable }

func (s *Device) Name() string { return deviceName }
func (s *Device) Value(ctx context.Context, _ interface{}, _ *cache.Cache) (interface{}, error) {
	device, ok := filterContext.FromDevice(ctx)
	if !ok {
		return nil, errors.New("device not found in context")
	}
	return device, nil
}
