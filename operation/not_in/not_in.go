package not_in

import (
	"context"
	"fmt"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/operation"
	"github.com/liyanbing/filter/utils"
	"github.com/liyanbing/filter/variable"
)

const Name = "nin"

var emptyElementErr = fmt.Errorf("[%s] expression must be greater than one element", Name)

func init() {
	operation.Register(&NotIn{})
}

type NotIn struct{}

func (s *NotIn) Name() string { return Name }
func (s *NotIn) PrepareValue(value interface{}) (interface{}, error) {
	targetValues := utils.ParseTargetArrayValue(value)
	if len(targetValues) == 0 {
		return nil, emptyElementErr
	}
	return targetValues, nil
}
func (s *NotIn) Run(ctx context.Context, v variable.Variable, value interface{}, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variable.GetValue(ctx, v, data, cache)
	if err != nil {
		return false, err
	}

	if targetValues, ok := value.([]interface{}); ok {
		for _, targetValue := range targetValues {
			if utils.ObjectCompare(variableValue, targetValue) == 0 {
				return false, nil
			}
		}
	}
	return true, nil
}
