package variables

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/liyanbing/filter/cache"
)

var getReg = regexp.MustCompile(`^get.(.+)`)

type Variable interface {
	Name() string
	Cacheable() bool
	Value(ctx context.Context, data interface{}, cache *cache.Cache) (interface{}, error)
}

type Builder interface {
	Name() string
	Build(string) Variable
}

type CalcFactorGetter interface {
	CalcValue(ctx context.Context, key string) (float64, error)
}

type FrequencyGetter interface {
	FrequencyValue(ctx context.Context, key string) interface{}
}

type Valuer interface {
	Value(ctx context.Context, key string) interface{}
}

// ============================== factory ==========================

var defaultFactory = &factory{
	builder: make(map[string]Builder),
}

type factory struct {
	builder map[string]Builder
	sync.Mutex
}

func (s *factory) Get(name string) (Variable, bool) {
	if builder, ok := s.builder[name]; ok {
		return builder.Build(name), true
	}

	segments := strings.Split(name, ".")
	if len(segments) <= 0 {
		return nil, false
	}

	if builder, ok := s.builder[segments[0]+"."]; ok {
		return builder.Build(name), true
	}
	return nil, false
}

func (s *factory) Register(builder Builder) {
	s.Lock()
	defer s.Unlock()
	name := builder.Name()
	if _, ok := s.builder[name]; ok {
		panic(fmt.Sprintf("%v variable already exists", name))
	}
	s.builder[name] = builder
}

func Register(builder Builder) {
	defaultFactory.Register(builder)
}

func Get(name string) (Variable, bool) {
	return defaultFactory.Get(name)
}

// ============================== Builder ==========================

type SimpleBuilder struct {
	name     string
	variable Variable
}

func NewSimpleVariable(variable Variable, opts ...Option) *SimpleBuilder {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	return &SimpleBuilder{
		name:     o.name,
		variable: variable,
	}
}

func (s *SimpleBuilder) Name() string {
	if s.name != "" {
		return s.name
	}
	return s.variable.Name()
}

func (s *SimpleBuilder) Build(_ string) Variable {
	return s.variable
}

type options struct {
	name string
}

type Option func(o *options)

func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

func GetValue(ctx context.Context, v Variable, data interface{}, cache *cache.Cache) (interface{}, error) {
	if v == nil {
		return nil, errors.New("empty variable")
	}

	if v.Cacheable() {
		if value, ok := cache.Get(v.Name()); ok {
			return value, nil
		}
	}

	value, err := v.Value(ctx, data, cache)
	if err != nil {
		return nil, err
	}

	if v.Cacheable() {
		cache.Set(v.Name(), value)
	}
	return value, nil
}
