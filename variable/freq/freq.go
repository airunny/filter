package freq

import (
	"context"
	"strings"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/variable"
)

const Name = "freq."

func init() {
	variable.Register(&freqBuilder{})
}

type freqBuilder struct{}

func (*freqBuilder) Name() string {
	return Name
}

func (*freqBuilder) Build(name string) variable.Variable {
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
	if getter, ok := data.(variable.FrequencyGetter); ok {
		freData := getter.FrequencyValue(ctx, s.key)
		return freData, nil
	}
	return 0, nil
}
