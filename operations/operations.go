package operations

import (
	"context"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/Liyanbing/filter/cache"
	"github.com/Liyanbing/filter/ip"
	"github.com/Liyanbing/filter/variables"
	"github.com/Liyanbing/filter/version"

	filterType "github.com/Liyanbing/filter/type"
)

type Operation interface {
	Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) bool
	PrepareValue(value interface{}) (interface{}, error)
}

type BaseOperationPrepareValue struct{}

func (s *BaseOperationPrepareValue) PrepareValue(value interface{}) (interface{}, error) {
	return value, nil
}

type factory struct {
	operations map[string]Operation
}

func (s *factory) Get(name string) Operation {
	if value, ok := s.operations[name]; ok {
		return value
	}

	return nil
}

func Register(name string, operation Operation) error {
	if _, ok := Factory.operations[name]; ok {
		return errors.New(name + " operation already exists")
	}

	Factory.operations[name] = operation
	return nil
}

// ----------------
var Factory *factory

func init() {
	Factory = &factory{
		operations: map[string]Operation{
			"=":       &Equal{},
			"!=":      &NotEqual{},
			"<>":      &NotEqual{},
			">":       &GreaterThan{},
			">=":      &GreaterThanEqual{},
			"<":       &LessThan{},
			"<=":      &LessThanEqual{},
			"~":       &Match{},
			"!~":      &NotMatch{},
			"~*":      &MatchAny{},
			"!~*":     &MatchNone{},
			"between": &Between{},
			"in":      &In{},
			"nin":     &NotIn{},
			"any":     &Any{},
			"has":     &Has{},
			"none":    &None{},
			"vgt":     &VersionGreaterThan{},
			"vgte":    &VersionGreaterThanOrEqual{},
			"vlt":     &VersionLessThan{},
			"vlte":    &VersionLessThanOrEqual{},
			"iir":     &InIPRange{},
			"niir":    &NotInIPRange{},
		},
	}
}

func GetVariableValue(ctx context.Context, v variables.Variable, data interface{}, cache *cache.Cache) interface{} {
	if v == nil {
		return ""
	}

	if v.Cacheable() {
		if value, ok := cache.Get(v.GetName()); ok {
			return value
		}
	}

	value := v.Value(ctx, data, cache)
	if v.Cacheable() {
		cache.Set(v.GetName(), value)
	}

	return value
}

func ParseTargetArrayValue(value interface{}) []interface{} {
	var target []interface{}

	switch filterType.GetMyType(value) {
	case filterType.STRING:
		targetValue := value.(string)
		values := strings.Split(targetValue, ",")
		for _, v := range values {
			target = append(target, strings.TrimSpace(v))
		}

	case filterType.ARRAY:
		target = value.([]interface{})
	}
	return target
}

//----------------------------------------------------------------------------------

// =
type Equal struct{ BaseOperationPrepareValue }

func (s *Equal) Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) bool {
	variableValue := GetVariableValue(ctx, variable, data, cache)

	return filterType.ObjectCompare(variableValue, value) == 0
}

// != (<>)
type NotEqual struct{ BaseOperationPrepareValue }

func (s *NotEqual) Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) bool {
	variableValue := GetVariableValue(ctx, variable, data, cache)

	return filterType.ObjectCompare(variableValue, value) != 0
}

// >
type GreaterThan struct{ BaseOperationPrepareValue }

func (s *GreaterThan) Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) bool {
	variableValue := GetVariableValue(ctx, variable, data, cache)

	return filterType.ObjectCompare(variableValue, value) == 1
}

// >=
type GreaterThanEqual struct{ BaseOperationPrepareValue }

func (s *GreaterThanEqual) Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) bool {
	variableValue := GetVariableValue(ctx, variable, data, cache)

	return filterType.ObjectCompare(variableValue, value) >= 0
}

// <
type LessThan struct{ BaseOperationPrepareValue }

func (s *LessThan) Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) bool {
	variableValue := GetVariableValue(ctx, variable, data, cache)

	return filterType.ObjectCompare(variableValue, value) == -1
}

// <=
type LessThanEqual struct{ BaseOperationPrepareValue }

func (s *LessThanEqual) Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) bool {
	variableValue := GetVariableValue(ctx, variable, data, cache)

	return filterType.ObjectCompare(variableValue, value) <= 0
}

// match
type Match struct{}

func (s *Match) Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) bool {
	variableValue := GetVariableValue(ctx, variable, data, cache)
	targetVariableValue, ok := variableValue.(string)
	if !ok {
		return false
	}

	if reg, ok := value.(*regexp.Regexp); ok {
		return reg.MatchString(targetVariableValue)

	} else if targetValue, ok := value.(string); ok {
		return strings.Contains(strings.ToLower(targetVariableValue), targetValue)

	}

	return false
}

func (s *Match) PrepareValue(value interface{}) (interface{}, error) {
	targetValue, ok := value.(string)
	if !ok {
		return nil, errors.New("invalid value")
	}

	if !(strings.HasPrefix(targetValue, "/") && strings.HasSuffix(targetValue, "/")) {
		return strings.ToLower(targetValue), nil
	}

	targetValue = strings.TrimSuffix(strings.TrimPrefix(targetValue, "/"), "/")
	if targetValue == "" {
		return nil, errors.New(fmt.Sprintf("[match] operation value is not a valid regexp expression[%s]", targetValue))
	}

	if reg, err := regexp.Compile("(?i)" + targetValue); err != nil {
		return nil, errors.New(fmt.Sprintf("[match] operation value is not a valid regexp expression[%s].err:%v", targetValue, err))
	} else {
		return reg, nil
	}
}

// not match
type NotMatch struct{}

func (s *NotMatch) Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) bool {
	variableValue := GetVariableValue(ctx, variable, data, cache)
	targetVariableValue, ok := variableValue.(string)
	if !ok {
		return false
	}

	if reg, ok := value.(*regexp.Regexp); ok {
		return !reg.MatchString(targetVariableValue)

	} else if targetValue, ok := value.(string); ok {
		return !strings.Contains(strings.ToLower(targetVariableValue), targetValue)

	}

	return false
}

func (s *NotMatch) PrepareValue(value interface{}) (interface{}, error) {
	targetValue, ok := value.(string)
	if !ok {
		return nil, errors.New("invalid value")
	}

	if !(strings.HasPrefix(targetValue, "/") && strings.HasSuffix(targetValue, "/")) {
		return strings.ToLower(targetValue), nil
	}

	targetValue = strings.TrimSuffix(strings.TrimPrefix(targetValue, "/"), "/")
	if targetValue == "" {
		return nil, errors.New(fmt.Sprintf("[not match] operation value is not a valid regexp expression[%s]", targetValue))
	}

	if reg, err := regexp.Compile("(?i)" + targetValue); err != nil {
		return nil, errors.New(fmt.Sprintf("[not match] operation value is not a valid regexp expression[%s].err:%v", targetValue, err))
	} else {
		return reg, nil
	}
}

// ~* (match any)
type MatchAny struct{}

func (s *MatchAny) Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) bool {
	variableValue := GetVariableValue(ctx, variable, data, cache)

	targetVariableValue, ok := variableValue.(string)
	if !ok {
		return false
	}

	targetValueElements, ok := value.([]interface{})
	if !ok {
		return false
	}

	for _, element := range targetValueElements {
		if reg, ok := element.(*regexp.Regexp); ok {
			if reg.MatchString(targetVariableValue) {
				return true
			}

		} else if targetValue, ok := element.(string); ok {
			if strings.Contains(strings.ToLower(targetVariableValue), targetValue) {
				return true
			}
		}
	}

	return false
}

func (s *MatchAny) PrepareValue(value interface{}) (interface{}, error) {
	targetValues := ParseTargetArrayValue(value)

	if len(targetValues) == 0 {
		return nil, errors.New(fmt.Sprintf("[ma] operation value must be a list"))
	}

	targetValueElements := make([]interface{}, 0, len(targetValues))
	for _, targetValue := range targetValues {
		targetValueStr, ok := targetValue.(string)
		if !ok {
			return nil, errors.New("[ma] operation value must be string")
		}

		if !(strings.HasPrefix(targetValueStr, "/") && strings.HasSuffix(targetValueStr, "/")) {
			targetValueElements = append(targetValueElements, strings.ToLower(targetValueStr))
			continue
		}

		targetValueStr = strings.TrimSuffix(strings.TrimPrefix(targetValueStr, "/"), "/")
		if targetValueStr == "" {
			return nil, errors.New(fmt.Sprintf("[ma] operation value is not a valid regexp expression[%s]", targetValue))
		}

		reg, err := regexp.Compile("(?i)" + targetValueStr)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("[ma] operation value is not a valid regexp expression[%s].err:%v", targetValue, err))
		}

		targetValueElements = append(targetValueElements, reg)
	}

	return targetValueElements, nil
}

// !~* (match none)
type MatchNone struct{}

func (s *MatchNone) Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) bool {
	variableValue := GetVariableValue(ctx, variable, data, cache)

	targetVariableValue, ok := variableValue.(string)
	if !ok {
		return false
	}

	targetValueElements, ok := value.([]interface{})
	if !ok {
		return false
	}

	for _, element := range targetValueElements {
		if reg, ok := element.(*regexp.Regexp); ok {
			if reg.MatchString(targetVariableValue) {
				return false
			}

		} else if targetValue, ok := element.(string); ok {
			if strings.Contains(strings.ToLower(targetVariableValue), targetValue) {
				return false
			}
		}
	}

	return true
}

func (s *MatchNone) PrepareValue(value interface{}) (interface{}, error) {
	targetValues := ParseTargetArrayValue(value)

	if len(targetValues) == 0 {
		return nil, errors.New(fmt.Sprintf("[ma] operation value must be a list"))
	}

	targetValueElements := make([]interface{}, 0, len(targetValues))
	for _, targetValue := range targetValues {
		targetValueStr, ok := targetValue.(string)
		if !ok {
			return nil, errors.New("[ma] operation value must be string")
		}

		if !(strings.HasPrefix(targetValueStr, "/") && strings.HasSuffix(targetValueStr, "/")) {
			targetValueElements = append(targetValueElements, strings.ToLower(targetValueStr))
			continue
		}

		targetValueStr = strings.TrimSuffix(strings.TrimPrefix(targetValueStr, "/"), "/")
		if targetValueStr == "" {
			return nil, errors.New(fmt.Sprintf("[ma] operation value is not a valid regexp expression[%s]", targetValue))
		}

		reg, err := regexp.Compile("(?i)" + targetValueStr)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("[ma] operation value is not a valid regexp expression[%s].err:%v", targetValue, err))
		}

		targetValueElements = append(targetValueElements, reg)
	}

	return targetValueElements, nil
}

// between
type Between struct{}

func (s *Between) Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) bool {
	variableValue := GetVariableValue(ctx, variable, data, cache)
	startAndEnd := value.([]interface{})

	return filterType.ObjectCompare(variableValue, startAndEnd[0]) >= 0 && filterType.ObjectCompare(variableValue, startAndEnd[1]) <= 0
}

func (s *Between) PrepareValue(value interface{}) (interface{}, error) {
	startAndEnd := ParseTargetArrayValue(value)

	if len(startAndEnd) != 2 {
		return nil, errors.New(fmt.Sprintf("[between] operation value must be greater than one element"))
	}

	return startAndEnd, nil
}

// in
type In struct{}

func (s *In) Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) bool {
	variableValue := GetVariableValue(ctx, variable, data, cache)

	if targetValues, ok := value.([]interface{}); ok {
		for _, targetValue := range targetValues {
			if filterType.ObjectCompare(variableValue, targetValue) == 0 {
				return true
			}
		}
		return false
	}

	return false
}

func (s *In) PrepareValue(value interface{}) (interface{}, error) {
	targetValues := ParseTargetArrayValue(value)

	if len(targetValues) == 0 {
		return nil, errors.New("[in] operation value must be greater than one element")
	}

	return targetValues, nil
}

// nin (not in)
type NotIn struct{}

func (s *NotIn) Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) bool {
	variableValue := GetVariableValue(ctx, variable, data, cache)

	if targetValues, ok := value.([]interface{}); ok {
		for _, targetValue := range targetValues {
			if filterType.ObjectCompare(variableValue, targetValue) == 0 {
				return false
			}
		}
		return true
	}

	return false
}

func (s *NotIn) PrepareValue(value interface{}) (interface{}, error) {
	targetValues := ParseTargetArrayValue(value)

	if len(targetValues) == 0 {
		return nil, errors.New("[nin] operation value must be greater than one element")
	}

	return targetValues, nil
}

// any
type Any struct{}

func (s *Any) Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) bool {
	variableValue := GetVariableValue(ctx, variable, data, cache)

	targetVariableValue, ok := variableValue.(string)
	if !ok {
		return false
	}

	variableValueElements := ParseTargetArrayValue(targetVariableValue)
	targetValueElements, ok := value.([]interface{})
	if !ok {
		return false
	}

	for _, targetValueElement := range targetValueElements {
		for _, variableValueElement := range variableValueElements {
			if filterType.ObjectCompare(targetValueElement, variableValueElement) == 0 {
				return true
			}
		}
	}

	return false
}

func (s *Any) PrepareValue(value interface{}) (interface{}, error) {
	targetValues := ParseTargetArrayValue(value)

	if len(targetValues) == 0 {
		return nil, errors.New("[any] operation value must be greater than one element")
	}

	return targetValues, nil
}

// has
type Has struct{}

func (s *Has) Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) bool {
	variableValue := GetVariableValue(ctx, variable, data, cache)

	targetVariableValue, ok := variableValue.(string)
	if !ok {
		return false
	}

	variableValueElements := ParseTargetArrayValue(targetVariableValue)
	targetValueElements, ok := value.([]interface{})
	if !ok {
		return false
	}

	for _, variableValueElement := range variableValueElements {
		has := false
		for _, valueElement := range targetValueElements {
			if filterType.ObjectCompare(variableValueElement, valueElement) == 0 {
				has = true
				break
			}
		}

		if !has {
			return false
		}
	}

	return true
}

func (s *Has) PrepareValue(value interface{}) (interface{}, error) {
	targetValues := ParseTargetArrayValue(value)

	if len(targetValues) == 0 {
		return nil, errors.New("[has] operation value must be greater than one element")
	}

	return targetValues, nil
}

// none
type None struct{}

func (s *None) Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) bool {
	variableValue := GetVariableValue(ctx, variable, data, cache)

	targetVariableValue, ok := variableValue.(string)
	if !ok {
		return false
	}

	variableValueElements := ParseTargetArrayValue(targetVariableValue)
	targetValueElements, ok := value.([]interface{})
	if !ok {
		return false
	}

	for _, variableValueElement := range variableValueElements {
		for _, valueElement := range targetValueElements {
			if filterType.ObjectCompare(variableValueElement, valueElement) == 0 {
				return false
			}
		}
	}

	return true
}

func (s *None) PrepareValue(value interface{}) (interface{}, error) {
	targetValues := ParseTargetArrayValue(value)

	if len(targetValues) == 0 {
		return nil, errors.New("[none] operation value must be greater than one element")
	}

	return targetValues, nil
}

// vgt (version greater than)
type VersionGreaterThan struct{ BaseOperationPrepareValue }

func (s *VersionGreaterThan) Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) bool {
	variableValue := GetVariableValue(ctx, variable, data, cache)

	if version.Compare(filterType.GetString(variableValue), filterType.GetString(value)) == 1 {
		return true
	}

	return false
}

// vgte (version greater than or equal)
type VersionGreaterThanOrEqual struct{ BaseOperationPrepareValue }

func (s *VersionGreaterThanOrEqual) Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) bool {
	variableValue := GetVariableValue(ctx, variable, data, cache)

	if version.Compare(filterType.GetString(variableValue), filterType.GetString(value)) != -1 {
		return true
	}

	return false
}

// vlt (version less than)
type VersionLessThan struct{ BaseOperationPrepareValue }

func (s *VersionLessThan) Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) bool {
	variableValue := GetVariableValue(ctx, variable, data, cache)

	if version.Compare(filterType.GetString(variableValue), filterType.GetString(value)) == -1 {
		return true
	}

	return false
}

// vlte (version less than or equal)
type VersionLessThanOrEqual struct{ BaseOperationPrepareValue }

func (s *VersionLessThanOrEqual) Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) bool {
	variableValue := GetVariableValue(ctx, variable, data, cache)

	if version.Compare(filterType.GetString(variableValue), filterType.GetString(value)) == -1 {
		return true
	}

	return false
}

// iir (in ip range)
type InIPRange struct{}

func (s *InIPRange) Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) bool {
	variableValue := GetVariableValue(ctx, variable, data, cache)

	targetVariableValue, ok := variableValue.(string)
	if !ok {
		return false
	}

	targetValue, ok := value.([]ip.Range)
	if !ok {
		return false
	}

	return ip.InRange(targetValue, net.ParseIP(targetVariableValue))
}

func (s *InIPRange) PrepareValue(value interface{}) (interface{}, error) {
	targetValue := ParseTargetArrayValue(value)

	if len(targetValue) == 0 {
		return nil, errors.New("[iir] operation value must be greater than one ip range")
	}

	ipRanges := make([]string, len(targetValue))
	for i, v := range targetValue {
		ipr, ok := v.(string)
		if !ok {
			return nil, errors.New("[iir] operation value must be a list of strings")
		}
		ipRanges[i] = ipr
	}

	return ip.Ranges(ipRanges...)
}

// niir (not in ip range)
type NotInIPRange struct{}

func (s *NotInIPRange) Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) bool {
	variableValue := GetVariableValue(ctx, variable, data, cache)

	targetVariableValue, ok := variableValue.(string)
	if !ok {
		return false
	}

	targetValue, ok := value.([]ip.Range)
	if !ok {
		return false
	}

	return !ip.InRange(targetValue, net.ParseIP(targetVariableValue))
}

func (s *NotInIPRange) PrepareValue(value interface{}) (interface{}, error) {
	targetValue := ParseTargetArrayValue(value)

	if len(targetValue) == 0 {
		return nil, errors.New("[niir] operation value must be greater than one ip range")
	}

	ipRanges := make([]string, len(targetValue))
	for i, v := range targetValue {
		ipr, ok := v.(string)
		if !ok {
			return nil, errors.New("[niir] operation value must be a list of strings")
		}
		ipRanges[i] = ipr
	}

	return ip.Ranges(ipRanges...)
}
