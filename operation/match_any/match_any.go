package match_any

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/operation"
	"github.com/liyanbing/filter/utils"
	"github.com/liyanbing/filter/variable"
)

const Name = "~*"

func init() {
	operation.Register(&MatchAny{})
}

type MatchAny struct{}

func (s *MatchAny) Name() string { return Name }
func (s *MatchAny) PrepareValue(value interface{}) (interface{}, error) {
	values := utils.ParseTargetArrayValue(value)
	if len(values) == 0 {
		return nil, fmt.Errorf("[%s] expression is invalid", value)
	}

	elements := make([]interface{}, 0, len(values))
	for _, v := range values {
		targetValueStr, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("[%s] expression must be string", Name)
		}

		if !(strings.HasPrefix(targetValueStr, "/") && strings.HasSuffix(targetValueStr, "/")) {
			elements = append(elements, targetValueStr)
			continue
		}

		targetValueStr = strings.TrimSuffix(strings.TrimPrefix(targetValueStr, "/"), "/")
		if targetValueStr == "" {
			return nil, fmt.Errorf("[%s] expression is not a valid regexp expression[%s]", Name, targetValueStr)
		}

		reg, err := regexp.Compile(targetValueStr)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("[%s] expression is not a valid regexp expression[%s]", Name, err))
		}
		elements = append(elements, reg)
	}
	return elements, nil
}

func (s *MatchAny) Run(ctx context.Context, v variable.Variable, value interface{}, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variable.GetValue(ctx, v, data, cache)
	if err != nil {
		return false, err
	}

	targetVariableValue, ok := variableValue.(string)
	if !ok {
		return false, fmt.Errorf("[%s] value must be string", Name)
	}

	elements, ok := value.([]interface{})
	if !ok {
		return false, fmt.Errorf("[%s] value must be array", Name)
	}

	for _, element := range elements {
		if reg, ok := element.(*regexp.Regexp); ok {
			if reg.MatchString(targetVariableValue) {
				return true, nil
			}
		} else if targetValue, ok := element.(string); ok {
			if strings.Contains(targetVariableValue, targetValue) {
				return true, nil
			}
		}
	}
	return false, nil
}
