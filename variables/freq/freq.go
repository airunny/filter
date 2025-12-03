package freq

import (
	"context"
	"strings"

	"github.com/airunny/filter/cache"
	"github.com/airunny/filter/variables"
)

const Name = "freq."

func init() {
	variables.Register(&freqBuilder{})
}

type freqBuilder struct{}

func (*freqBuilder) Name() string {
	return Name
}

func (*freqBuilder) Build(name string) variables.Variable {
	key := strings.TrimPrefix(name, Name)
	if key == "" {
		return nil
	}
	return &Freq{
		name: name,
		key:  key,
	}
}

// Freq 频次控制
type Freq struct {
	name string
	key  string
}

func (s *Freq) Name() string    { return s.name }
func (s *Freq) Cacheable() bool { return false }
func (s *Freq) Value(ctx context.Context, data interface{}, _ *cache.Cache) (interface{}, error) {
	if getter, ok := data.(variables.Frequency); ok {
		return getter.FrequencyValue(ctx, s.key)
	}
	return 0, nil
}
