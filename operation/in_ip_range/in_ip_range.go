package in_ip_range

import (
	"context"
	"fmt"
	"net"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/ip"
	"github.com/liyanbing/filter/operation"
	"github.com/liyanbing/filter/utils"
	"github.com/liyanbing/filter/variable"
)

const Name = "iir"

func init() {
	operation.Register(&InIPRange{})
}

type InIPRange struct{}

func (s *InIPRange) Name() string { return Name }
func (s *InIPRange) PrepareValue(value interface{}) (interface{}, error) {
	targetValue := utils.ParseTargetArrayValue(value)
	if len(targetValue) == 0 {
		return nil, fmt.Errorf("[%s] value must be a list of string", Name)
	}

	ipRanges := make([]string, len(targetValue))
	for i, v := range targetValue {
		ipr, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("[%s] value must be string", Name)
		}
		ipRanges[i] = ipr
	}
	return ip.Ranges(ipRanges...)
}
func (s *InIPRange) Run(ctx context.Context, v variable.Variable, value interface{}, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variable.GetValue(ctx, v, data, cache)
	if err != nil {
		return false, err
	}

	targetVariableValue, ok := variableValue.(string)
	if !ok {
		return false, fmt.Errorf("[%s] variable value must be string", Name)
	}

	targetValue, ok := value.([]ip.Range)
	if !ok {
		return false, fmt.Errorf("[%s] value must be a list of string", Name)
	}
	return ip.InRange(targetValue, net.ParseIP(targetVariableValue)), nil
}
