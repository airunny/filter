package filter

import (
	"context"
	"errors"
	"sort"

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

func BuildFilter(ctx context.Context, filterData []interface{}) (*Filter, error) {
	if len(filterData) < 2 {
		return nil, errors.New("filterData struct mast contain at least 2 items")
	}

	filterCount := len(filterData)
	filterCondition, err := condition.BuildCondition(ctx, filterData[:filterCount-1], condition.LOGIC_AND)
	if err != nil {
		return nil, err
	}

	filterExecutor, err := executor.BuildExecutor(ctx, filterData[filterCount-1:])
	if err != nil {
		return nil, err
	}

	return &Filter{
		condition: filterCondition,
		executor:  filterExecutor,
	}, nil
}

// --------------
type filterPack struct {
	Filter   *Filter
	Id       string
	Weight   int64
	Priority int64
}

func (s *filterPack) GetWeight() int64 {
	return s.Weight
}

func (s *filterPack) GetPriority() int64 {
	return s.Priority
}

// ----------------
type priorityBoundary struct {
	NextPriorityStartIndex int
	TotalWeights           int64
}

type GroupFilter struct {
	filters          []*filterPack
	priorityBoundary []priorityBoundary
	batchMode        bool
	totalWeights     int64
}

func NewEmptyFilterGroup(batchMode bool) *GroupFilter {
	return &GroupFilter{
		batchMode: batchMode,
	}
}

func NewGroupFilterWithConfig(ctx context.Context, cnf *Config) (*GroupFilter, error) {
	group := NewEmptyFilterGroup(false)
	for filterID, filterCnf := range cnf.Filters {
		singleFilter, err := BuildFilter(ctx, filterCnf.FilterData)
		if err != nil {
			return nil, err
		}

		group.Add(singleFilter, filterID, filterCnf.Priority, filterCnf.Weight)
	}

	return group, nil
}

func (s *GroupFilter) Run(ctx context.Context, data interface{}, cache *cache.Cache) (successNumber int, filterID string) {
	if s.totalWeights > 0 {
		filters := make([]utils.IWeight, 0, len(s.filters))
		for _, filter := range s.filters {
			filters = append(filters, filter)
		}

		lastBoundary := 0
		for _, boundary := range s.priorityBoundary {
			if boundary.TotalWeights != 0 {
				utils.ShuffleByWeight(filters[lastBoundary:boundary.NextPriorityStartIndex], boundary.TotalWeights)
			}
			lastBoundary = boundary.NextPriorityStartIndex
		}

		s.filters = s.filters[:0]
		for _, filter := range filters {
			s.filters = append(s.filters, filter.(*filterPack))
		}
	}

	for _, filter := range s.filters {
		if filter.Filter.Run(ctx, data, cache) {
			filterID = filter.Id
			successNumber++
			if !s.batchMode {
				break
			}
		}
	}
	return
}

func (s *GroupFilter) Add(filter *Filter, filterId string, priority int64, weight int64) {
	s.filters = append(s.filters, &filterPack{
		Filter:   filter,
		Id:       filterId,
		Weight:   weight,
		Priority: priority,
	})
	s.totalWeights += weight

	sort.Slice(s.filters, func(i, j int) bool {
		return s.filters[i].GetPriority() > s.filters[j].GetPriority()
	})

	s.locatePriorityBoundary()
}

func (s *GroupFilter) locatePriorityBoundary() {
	s.priorityBoundary = s.priorityBoundary[:0]
	lastPriority := int64(0)
	totalWeight := int64(0)

	for index, filter := range s.filters {
		if index != 0 {
			if filter.Priority != lastPriority {
				s.priorityBoundary = append(s.priorityBoundary, priorityBoundary{
					NextPriorityStartIndex: index,
					TotalWeights:           totalWeight,
				})
				totalWeight = 0
			}
		}
		totalWeight += filter.GetWeight()
		lastPriority = filter.Priority
	}

	s.priorityBoundary = append(s.priorityBoundary, priorityBoundary{
		NextPriorityStartIndex: len(s.filters),
		TotalWeights:           totalWeight,
	})
}
