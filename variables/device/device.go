package device

import (
	"context"
	"errors"

	"github.com/airunny/filter/cache"
	filterContext "github.com/airunny/filter/context"
	"github.com/airunny/filter/variables"
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
