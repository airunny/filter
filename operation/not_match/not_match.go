package not_match

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/operation"
	"github.com/liyanbing/filter/variable"
)

const Name = "!~"

func init() {
	operation.Register(&NotMatch{})
}

type NotMatch struct{}

func (s *NotMatch) Name() string { return Name }
func (s *NotMatch) PrepareValue(value interface{}) (interface{}, error) {
	targetValue, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("[%s] expression must be string", Name)
	}

	if !(strings.HasPrefix(targetValue, "/") && strings.HasSuffix(targetValue, "/")) {
		return targetValue, nil
	}

	targetValue = strings.TrimSuffix(strings.TrimPrefix(targetValue, "/"), "/")
	if targetValue == "" {
		return nil, fmt.Errorf("[%s] expression is not a valid regexp expression[%s]", Name, value)
	}

	reg, err := regexp.Compile(targetValue)
	if err != nil {
		return nil, fmt.Errorf("[%s] expression is not a valid regexp expression[%s]", Name, value)
	}
	return reg, nil
}

func (s *NotMatch) Run(ctx context.Context, v variable.Variable, value interface{}, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variable.GetValue(ctx, v, data, cache)
	if err != nil {
		return false, err
	}

	targetVariableValue, ok := variableValue.(string)
	if !ok {
		return false, fmt.Errorf("[%s] value must be string", Name)
	}

	if reg, ok := value.(*regexp.Regexp); ok {
		return !reg.MatchString(targetVariableValue), nil
	} else if targetValue, ok := value.(string); ok {
		return !strings.Contains(targetVariableValue, targetValue), nil
	} else {
		return false, fmt.Errorf("[%s] operation value must be string", Name)
	}
}
