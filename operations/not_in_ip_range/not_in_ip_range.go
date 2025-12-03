package not_in_ip_range

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/airunny/filter/cache"
	"github.com/airunny/filter/operations"
	"github.com/airunny/filter/utils"
	"github.com/airunny/filter/variables"
)

const Name = "niir"

var (
	ErrInvalidOperationValue        = fmt.Errorf("[%s] operation value must be a list of string", Name)
	ErrInvalidOperationValueElement = fmt.Errorf("[%s] variable value must be string", Name)
	ErrEmptyOperationValueElement   = fmt.Errorf("[%s] variable value must be not empty string", Name)
	ErrInvalidVariableValue         = fmt.Errorf("[%s] variable value must be string", Name)
)

func init() {
	operations.Register(&NotInIPRange{})
}

type NotInIPRange struct{}

func (s *NotInIPRange) Name() string { return Name }
func (s *NotInIPRange) PrepareValue(value interface{}) (interface{}, error) {
	targetValue := utils.ParseTargetArrayValue(value)
	if len(targetValue) == 0 {
		return nil, ErrInvalidOperationValue
	}

	ipRanges := make([]string, 0, len(targetValue))
	for _, v := range targetValue {
		ipr, ok := v.(string)
		if !ok {
			return nil, ErrInvalidOperationValueElement
		}

		ipr = strings.TrimSpace(ipr)
		if ipr == "" {
			return nil, ErrEmptyOperationValueElement
		}
		ipRanges = append(ipRanges, ipr)
	}
	return utils.IPRanges(ipRanges...)
}

func (s *NotInIPRange) Run(ctx context.Context, variable variables.Variable, operationValue, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variables.GetValue(ctx, variable, data, cache)
	if err != nil {
		return false, err
	}

	targetVariableValue, ok := variableValue.(string)
	if !ok {
		return false, ErrInvalidVariableValue
	}

	targetValue, ok := operationValue.([]utils.IPRange)
	if !ok {
		return false, ErrInvalidOperationValue
	}
	return !utils.InIPRange(targetValue, net.ParseIP(targetVariableValue)), nil
}
