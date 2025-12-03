package not_match

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/operations"
	"github.com/liyanbing/filter/variables"
)

const Name = "!~"

var (
	ErrInvalidOperationValue = fmt.Errorf("[%s] operation value must be string", Name)
	ErrInvalidVariableValue  = fmt.Errorf("[%s] variable value must be string", Name)
)

func init() {
	operations.Register(&NotMatch{})
}

type NotMatch struct{}

func (s *NotMatch) Name() string { return Name }
func (s *NotMatch) PrepareValue(value interface{}) (interface{}, error) {
	targetValue, ok := value.(string)
	if !ok {
		return nil, ErrInvalidOperationValue
	}

	if !(strings.HasPrefix(targetValue, "/") && strings.HasSuffix(targetValue, "/")) {
		return targetValue, nil
	}

	targetValue = strings.TrimSuffix(strings.TrimPrefix(targetValue, "/"), "/")
	if targetValue == "" {
		return nil, fmt.Errorf("[%s] operation value is not a valid regexp [%s]", Name, value)
	}

	reg, err := regexp.Compile(targetValue)
	if err != nil {
		return nil, fmt.Errorf("[%s] operation value is not a valid regexp [%s]", Name, value)
	}
	return reg, nil
}

func (s *NotMatch) Run(ctx context.Context, variable variables.Variable, operationValue, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variables.GetValue(ctx, variable, data, cache)
	if err != nil {
		return false, err
	}

	targetVariableValue, ok := variableValue.(string)
	if !ok {
		return false, ErrInvalidVariableValue
	}

	if reg, ok := operationValue.(*regexp.Regexp); ok {
		return !reg.MatchString(targetVariableValue), nil
	} else if targetValue, ok := operationValue.(string); ok {
		return !strings.Contains(targetVariableValue, targetValue), nil
	} else {
		return false, ErrInvalidOperationValue
	}
}
