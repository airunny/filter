package between

import (
	"context"
	"fmt"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/operations"
	"github.com/liyanbing/filter/utils"
	"github.com/liyanbing/filter/variables"
)

const Name = "between"

var ErrInvalidOperationValue = fmt.Errorf("[%s] operation value must have two element", Name)

func init() {
	operations.Register(&Between{})
}

type Between struct{}

func (s *Between) Name() string { return Name }
func (s *Between) PrepareValue(value interface{}) (interface{}, error) {
	elements := utils.ParseTargetArrayValue(value)
	if len(elements) != 2 {
		return nil, ErrInvalidOperationValue
	}
	return elements, nil
}

func (s *Between) Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) (bool, error) {
	variableValue, err := variables.GetValue(ctx, variable, data, cache)
	if err != nil {
		return false, err
	}
	startAndEnd, ok := value.([]interface{})
	if !ok {
		return false, ErrInvalidOperationValue
	}
	return utils.ObjectCompare(variableValue, startAndEnd[0]) >= 0 && utils.ObjectCompare(variableValue, startAndEnd[1]) <= 0, nil
}
