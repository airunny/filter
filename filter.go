package filter

import (
	"context"
	"errors"
	"sort"
	"sync"

	"github.com/Liyanbing/filter/cache"
	"github.com/Liyanbing/filter/condition"
	"github.com/Liyanbing/filter/executor"
	"github.com/Liyanbing/filter/utils"
)

type Filter struct {
	condition condition.Condition
	executor  executor.Executor
}

func (s *Filter) Run(ctx context.Context, data interface{}, cache *cache.Cache) bool {
	if !s.condition.IsConditionOk(ctx, data, cache) {
		return false
	}

	s.executor.Execute(ctx, data)
	return true
}

func (s *Filter) SingleRun(data interface{}, ctx context.Context, cache *cache.Cache) bool {
	return s.Run(ctx, data, cache)
}

// ----------------
type GroupFilter struct {
	Lock             sync.Mutex
	filters          []*filter2W
	priorityBoundary []boundary
	batchMode        bool
	totalWeights     int64
}

type boundary struct {
	Index   int
	Weights int64
}

type filter2W struct {
	Filter   *Filter
	Id       string
	Weight   int64
	Priority int64
}

func NewEmptyFilterGroup(batchMode bool) *GroupFilter {
	return &GroupFilter{
		batchMode: batchMode,
	}
}

func (s *filter2W) GetWeight() int64 {
	return s.Weight
}

func (s *filter2W) GetPriority() int64 {
	return s.Priority
}

func (s *GroupFilter) Run(ctx context.Context, data interface{}, cache *cache.Cache) (successCount int, successID string) {
	if s.totalWeights == 0 {
		for _, filter := range s.filters {
			if filter.Filter.Run(ctx, data, cache) {
				successID = filter.Id
				successCount++
				if !s.batchMode {
					break
				}
			}
		}
	} else {
		filters := make([]utils.IWeight, 0, len(s.filters))
		for _, filter := range s.filters {
			filters = append(filters, filter)
		}

		lastBoundary := 0
		for _, boundary := range s.priorityBoundary {
			if boundary.Weights != 0 {
				utils.ShuffleByWeight(filters[lastBoundary:boundary.Index], boundary.Weights)
			}
			lastBoundary = boundary.Index
		}

		for _, filter := range filters {
			filter := filter.(*filter2W)
			if filter.Filter.Run(ctx, data, cache) {
				successID = filter.Id
				successCount++
				if !s.batchMode {
					break
				}
			}
		}
	}
	return
}

func (s *GroupFilter) add(filter *Filter, filterId string, priority int64, weight int64) {
	s.filters = append(s.filters, &filter2W{
		Id:       filterId,
		Filter:   filter,
		Weight:   weight,
		Priority: priority,
	})
	s.totalWeights += weight
}

func (s *GroupFilter) Add(filter *Filter, filterId string, p int64, w int64) {
	s.add(filter, filterId, p, w)

	sort.Slice(s.filters, func(i, j int) bool {
		return s.filters[i].GetPriority() > s.filters[j].GetPriority()
	})

	s.locatePriorityBoundary()
}

func (s *GroupFilter) locatePriorityBoundary() {
	s.priorityBoundary = s.priorityBoundary[:0]
	lastP := int64(0)
	totalWeight := int64(0)

	for index, filter := range s.filters {
		if index != 0 {
			if filter.Priority != lastP {
				s.priorityBoundary = append(s.priorityBoundary, boundary{
					Index:   index,
					Weights: totalWeight,
				})
				totalWeight = 0
			}
		}
		totalWeight += filter.GetWeight()
		lastP = filter.Priority
	}
	s.priorityBoundary = append(s.priorityBoundary, boundary{len(s.filters), totalWeight}) // 最后加上最后一个元素
}

func BuildFilter(filter []interface{}) (*Filter, error) {
	if len(filter) < 2 {
		return nil, errors.New("filter struct mast contain at least 2 items")
	}

	filterCount := len(filter)
	con, err := condition.BuildCondition(filter[:filterCount-1], condition.LOGIC_AND)
	if err != nil {
		return nil, err
	}

	exe, err := executor.BuildExecutor(filter[filterCount-1:])
	if err != nil {
		return nil, err
	}

	return &Filter{
		condition: con,
		executor:  exe,
	}, nil
}
