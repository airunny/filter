package filter

import (
	"context"
	"errors"
	"math/rand"
	"sort"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/condition"
	"github.com/liyanbing/filter/executor"
)

type Filter struct {
	id        string
	weight    int64
	priority  int64
	condition condition.Condition
	executor  executor.Executor
}

func (s *Filter) Weight() int64 {
	return s.weight
}

func (s *Filter) Priority() int64 {
	return s.priority
}

func (s *Filter) Run(ctx context.Context, data interface{}, cache *cache.Cache) (bool, error) {
	ok, err := s.condition.IsConditionOk(ctx, data, cache)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, nil
	}

	err = s.executor.Execute(ctx, data)
	if err != nil {
		return false, err
	}
	return true, nil
}

func BuildFilter(ctx context.Context, id string, weight, priority int64, filterData []interface{}) (*Filter, error) {
	if len(filterData) < 2 {
		return nil, errors.New("filter must contain at least two items")
	}

	filterCount := len(filterData)
	filterCondition, err := condition.BuildCondition(ctx, filterData[:filterCount-1], condition.LogicAnd)
	if err != nil {
		return nil, err
	}

	filterExecutor, err := executor.BuildExecutor(ctx, filterData[filterCount-1:])
	if err != nil {
		return nil, err
	}

	return &Filter{
		id:        id,
		weight:    weight,
		priority:  priority,
		condition: filterCondition,
		executor:  filterExecutor,
	}, nil
}

type priorityBoundary struct {
	nextIndex int
	weight    int64
}

type GroupFilter struct {
	filters          []*Filter
	priorityBoundary []priorityBoundary
	batch            bool
	weight           int64
}

func NewGroupFilterWithConfig(ctx context.Context, cnf *Config) (*GroupFilter, error) {
	group := &GroupFilter{
		filters: make([]*Filter, 0, len(cnf.Filters)),
		batch:   cnf.Batch,
	}

	for id, filter := range cnf.Filters {
		singleFilter, err := BuildFilter(ctx, id, filter.Weight, filter.Priority, filter.Filter)
		if err != nil {
			return nil, err
		}
		group.Add(singleFilter)
	}
	return group, nil
}

func (s *GroupFilter) Run(ctx context.Context, data interface{}, cache *cache.Cache) (successNumber int, filterIds []string, err error) {
	filters := s.filters
	if s.weight > 0 {
		lastBoundary := 0
		filters = make([]*Filter, 0, len(s.filters))
		for _, filter := range filters {
			filters = append(filters, filter)
		}

		for _, boundary := range s.priorityBoundary {
			if boundary.weight != 0 {
				shuffleByWeight(filters[lastBoundary:boundary.nextIndex], boundary.weight)
			}
			lastBoundary = boundary.nextIndex
		}
	}

	for _, filter := range filters {
		var ok bool
		ok, err = filter.Run(ctx, data, cache)
		if err != nil {
			return
		}

		if !ok {
			continue
		}

		filterIds = append(filterIds, filter.id)
		successNumber++
		if !s.batch {
			break
		}
	}
	return
}

func (s *GroupFilter) Add(filter *Filter) {
	s.filters = append(s.filters, filter)
	s.weight += filter.weight

	sort.Slice(s.filters, func(i, j int) bool {
		return s.filters[i].Priority() > s.filters[j].Priority()
	})
	s.locatePriorityBoundary()
}

func (s *GroupFilter) locatePriorityBoundary() {
	s.priorityBoundary = s.priorityBoundary[:0]
	lastPriority := int64(0)
	totalWeight := int64(0)

	for index, filter := range s.filters {
		if index != 0 {
			if filter.Priority() != lastPriority {
				s.priorityBoundary = append(s.priorityBoundary, priorityBoundary{
					nextIndex: index,
					weight:    totalWeight,
				})
				totalWeight = 0
			}
		}
		totalWeight += filter.Weight()
		lastPriority = filter.Priority()
	}

	s.priorityBoundary = append(s.priorityBoundary, priorityBoundary{
		nextIndex: len(s.filters),
		weight:    totalWeight,
	})
}

func filterWeight(filters []*Filter) int64 {
	total := int64(0)
	for _, filter := range filters {
		total += filter.Weight()
	}
	return total
}

func pickByWeight(filters []*Filter, totalWeight int64) int {
	if totalWeight == 0 {
		totalWeight = filterWeight(filters)
	}

	var (
		choose = rand.Int63n(totalWeight) + 1
		line   = int64(0)
	)

	for i, filter := range filters {
		line += filter.Weight()
		if choose <= line {
			return i
		}
	}
	return 0
}

func shuffleByWeight(filters []*Filter, totalWeight int64) {
	if len(filters) == 0 || len(filters) == 1 {
		return
	}

	if totalWeight == 0 {
		totalWeight = filterWeight(filters)
	}

	for curIndex := 0; curIndex < len(filters); curIndex++ {
		chooseIndex := curIndex + pickByWeight(filters[curIndex:], totalWeight)
		filters[chooseIndex], filters[curIndex] = filters[curIndex], filters[chooseIndex]
		totalWeight -= filters[curIndex].Weight()
	}
}
