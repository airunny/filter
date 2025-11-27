package match_none

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/operations"
	"github.com/liyanbing/filter/utils"
	"github.com/liyanbing/filter/variables"
)

const Name = "!~*"

var (
	ErrInvalidOperationValue        = fmt.Errorf("[%s] operation value must be string", Name)
	ErrInvalidOperationElementValue = fmt.Errorf("[%s] operation value item must be string", Name)
	ErrEmptyOperationValue          = fmt.Errorf("[%s] operation value can not be empty", Name)
	ErrInvalidVariableValue         = fmt.Errorf("[%s] variable value must be string", Name)
)

func init() {
	operations.Register(&MatchNone{})
}

type MatchNone struct{}

func (s *MatchNone) Name() string { return Name }
func (s *MatchNone) PrepareValue(value interface{}) (interface{}, error) {
	values := utils.ParseTargetArrayValue(value)
	if len(values) == 0 {
		return nil, ErrInvalidOperationValue
	}

	elements := make([]interface{}, 0, len(values))
	for _, v := range values {
		targetValueStr, ok := v.(string)
		if !ok {
			return nil, ErrInvalidOperationElementValue
		}

		if !(strings.HasPrefix(targetValueStr, "/") && strings.HasSuffix(targetValueStr, "/")) {
			elements = append(elements, targetValueStr)
			continue
		}

		targetValueStr = strings.TrimSuffix(strings.TrimPrefix(targetValueStr, "/"), "/")
		if targetValueStr == "" {
			return nil, ErrEmptyOperationValue
		}

		reg, err := regexp.Compile(targetValueStr)
		if err != nil {
			return nil, fmt.Errorf("[%s] operation value invalid regexp [%s]", Name, targetValueStr)
		}
		elements = append(elements, reg)
	}
	return elements, nil
}

func (s *MatchNone) Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) (bool, error) {
	elements, ok := value.([]interface{})
	if !ok {
		return false, ErrInvalidOperationValue
	}

	variableValue, err := variables.GetValue(ctx, variable, data, cache)
	if err != nil {
		return false, err
	}

	targetVariableValue, ok := variableValue.(string)
	if !ok {
		return false, ErrInvalidVariableValue
	}

	for _, element := range elements {
		if reg, ok := element.(*regexp.Regexp); ok {
			if reg.MatchString(targetVariableValue) {
				return false, nil
			}
		} else if targetValue, ok := element.(string); ok {
			if strings.Contains(targetVariableValue, targetValue) {
				return false, nil
			}
		}
	}
	return true, nil
}
