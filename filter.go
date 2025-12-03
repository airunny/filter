package filter

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"sort"
	"strings"
	"sync/atomic"

	"github.com/airunny/filter/cache"
	"github.com/airunny/filter/condition"
	"github.com/airunny/filter/executor"
)

type Reporter interface {
	Report(ctx context.Context, data interface{}, filterIds []string)
}

type ReportFunc func(ctx context.Context, data interface{}, filterIds []string)

func (f ReportFunc) Report(ctx context.Context, data interface{}, filterIds []string) {
	f(ctx, data, filterIds)
}

type Config struct {
	Filters []struct {
		Id       string        `json:"id"`
		Weight   int64         `json:"weight"`
		Priority int64         `json:"priority"`
		Filter   []interface{} `json:"Filter"`
	} `json:"filters"`
	Batch bool `json:"batch"`
}

type Filter struct {
	batch    atomic.Value
	reporter Reporter
}

func NewFilter(ctx context.Context, jsonStr string, reporter Reporter) (*Filter, error) {
	var cnf Config
	err := json.NewDecoder(strings.NewReader(jsonStr)).Decode(&cnf)
	if err != nil {
		return nil, err
	}

	batch, err := buildBatchFilter(ctx, &cnf)
	if err != nil {
		return nil, err
	}

	batchValue := atomic.Value{}
	batchValue.Store(batch)
	return &Filter{
		batch:    batchValue,
		reporter: reporter,
	}, nil
}

func (s *Filter) Execute(ctx context.Context, data interface{}) (interface{}, error) {
	batch, ok := s.batch.Load().(*batchFilter)
	if !ok {
		return nil, errors.New("invalid Filter")
	}

	if data == nil {
		data = make(map[string]interface{})
	}

	_, filterIds, err := batch.Run(ctx, data, cache.NewCache())
	if err != nil {
		return nil, err
	}

	if s.reporter != nil {
		s.reporter.Report(ctx, data, filterIds)
	}
	return data, nil
}

func (s *Filter) Refresh(ctx context.Context, jsonStr string) error {
	var cnf Config
	err := json.NewDecoder(strings.NewReader(jsonStr)).Decode(&cnf)
	if err != nil {
		return err
	}

	batch, err := buildBatchFilter(ctx, &cnf)
	if err != nil {
		return err
	}

	s.batch.Store(batch)
	return nil
}

type singleFilter struct {
	id        string
	weight    int64
	priority  int64
	condition condition.Condition
	executor  executor.Executor
}

func (s *singleFilter) Weight() int64 {
	return s.weight
}

func (s *singleFilter) Priority() int64 {
	return s.priority
}

func (s *singleFilter) Run(ctx context.Context, data interface{}, cache *cache.Cache) (bool, error) {
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

type priorityBoundary struct {
	nextIndex int
	weight    int64
}

type batchFilter struct {
	filters    []*singleFilter
	priorities []priorityBoundary
	batch      bool
	weight     int64
}

func buildBatchFilter(ctx context.Context, cnf *Config) (*batchFilter, error) {
	batch := &batchFilter{
		filters: make([]*singleFilter, 0, len(cnf.Filters)),
		batch:   cnf.Batch,
	}

	for _, filter := range cnf.Filters {
		single, err := buildSingleFilter(ctx, filter.Id, filter.Weight, filter.Priority, filter.Filter)
		if err != nil {
			return nil, err
		}
		batch.Add(single)
	}
	return batch, nil
}

func (s *batchFilter) Run(ctx context.Context, data interface{}, cache *cache.Cache) (successNumber int, filterIds []string, err error) {
	filters := s.filters
	if s.weight > 0 {
		lastBoundary := 0
		for _, boundary := range s.priorities {
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

func (s *batchFilter) Add(filter *singleFilter) {
	s.filters = append(s.filters, filter)
	s.weight += filter.weight

	sort.Slice(s.filters, func(i, j int) bool {
		return s.filters[i].Priority() < s.filters[j].Priority()
	})
	s.locatePriority()
}

func (s *batchFilter) locatePriority() {
	s.priorities = s.priorities[:0]
	lastPriority := int64(0)
	totalWeight := int64(0)

	for index, filter := range s.filters {
		if index != 0 {
			if filter.Priority() != lastPriority {
				s.priorities = append(s.priorities, priorityBoundary{
					nextIndex: index,
					weight:    totalWeight,
				})
				totalWeight = 0
			}
		}
		totalWeight += filter.Weight()
		lastPriority = filter.Priority()
	}
	s.priorities = append(s.priorities, priorityBoundary{
		nextIndex: len(s.filters),
		weight:    totalWeight,
	})
}

func buildSingleFilter(ctx context.Context, id string, weight, priority int64, filterData []interface{}) (*singleFilter, error) {
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

	return &singleFilter{
		id:        id,
		weight:    weight,
		priority:  priority,
		condition: filterCondition,
		executor:  filterExecutor,
	}, nil
}

func filterWeight(filters []*singleFilter) int64 {
	total := int64(0)
	for _, filter := range filters {
		total += filter.Weight()
	}
	return total
}

func pickByWeight(filters []*singleFilter, totalWeight int64) int {
	var (
		choose = rand.Int63n(totalWeight) + 1
		line   = int64(0)
	)

	for i, f := range filters {
		line += f.Weight()
		if choose <= line {
			return i
		}
	}
	return 0
}

func shuffleByWeight(filters []*singleFilter, totalWeight int64) {
	if len(filters) == 0 || len(filters) == 1 {
		return
	}

	for curIndex := 0; curIndex < len(filters); curIndex++ {
		chooseIndex := curIndex + pickByWeight(filters[curIndex:], totalWeight)
		totalWeight -= filters[chooseIndex].Weight()
		if chooseIndex == curIndex {
			continue
		}
		filters[chooseIndex], filters[curIndex] = filters[curIndex], filters[chooseIndex]
	}
}
