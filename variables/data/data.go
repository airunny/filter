package data

import (
	"context"
	"fmt"
	"strings"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/utils"
	"github.com/liyanbing/filter/variables"
)

const Name = "data."

func init() {
	variables.Register(&dataBuilder{})
}

type dataBuilder struct{}

func (*dataBuilder) Name() string {
	return Name
}

func (*dataBuilder) Build(name string) variables.Variable {
	key := strings.TrimPrefix(name, Name)
	if key == "" {
		return nil
	}
	return &Data{
		name: name,
		key:  key,
	}
}

// Data 获取传递的data中的值
type Data struct {
	name string
	key  string
}

func (s *Data) Name() string    { return s.name }
func (s *Data) Cacheable() bool { return false }
func (s *Data) Value(ctx context.Context, data interface{}, _ *cache.Cache) (interface{}, error) {
	value, ok := utils.GetObjectValueByKey(ctx, data, s.key)
	if !ok {
		return nil, fmt.Errorf("%s not found in data", s.name)
	}
	return value, nil
}
