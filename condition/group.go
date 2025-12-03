package condition

import (
	"context"
	"errors"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/types"
)

func NewGroup(logic Logic) *Group {
	return &Group{
		logic:      logic,
		conditions: make([]Condition, 0, 2),
	}
}

type Group struct {
	logic      Logic
	conditions []Condition
}

func (s *Group) Add(condition Condition) {
	s.conditions = append(s.conditions, condition)
}

func (s *Group) IsConditionOk(ctx context.Context, data interface{}, cache *cache.Cache) (bool, error) {
	result := true
	for _, condition := range s.conditions {
		ok, err := condition.IsConditionOk(ctx, data, cache)
		if err != nil {
			return false, err
		}

		if ok {
			if s.logic == LogicAnd {
				continue
			}

			if s.logic == LogicOr {
				result = true
				break
			}

			if s.logic == LogicNot {
				result = false
				break
			}
		} else {
			if s.logic == LogicAnd {
				result = false
				break
			}

			if s.logic == LogicOr {
				result = false
				continue
			}

			if s.logic == LogicNot {
				continue
			}
		}
	}
	return result, nil
}

func BuildGroup(ctx context.Context, items []interface{}, logic Logic) (Condition, error) {
	group := NewGroup(logic)
	for _, item := range items {
		if !types.IsArray(item) {
			return nil, errors.New("condition item is not array")
		}

		subCondition, err := BuildCondition(ctx, item.([]interface{}), LogicAnd)
		if err != nil {
			return nil, err
		}
		group.Add(subCondition)
	}
	return group, nil
}
