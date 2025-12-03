package device

import (
	"context"
	"errors"

	"github.com/liyanbing/filter/cache"
	filterContext "github.com/liyanbing/filter/context"
	"github.com/liyanbing/filter/variables"
)

const Name = "device"

func init() {
	variables.Register(variables.NewSimpleVariable(&Device{}))
}

// Device 设备ID
type Device struct{}

func (s *Device) Name() string    { return Name }
func (s *Device) Cacheable() bool { return true }
func (s *Device) Value(ctx context.Context, _ interface{}, _ *cache.Cache) (interface{}, error) {
	device, ok := filterContext.FromDevice(ctx)
	if !ok {
		return nil, errors.New("device not found in context")
	}
	return device, nil
}
