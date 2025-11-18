package freq

import (
	"context"
	"strings"

	"github.com/liyanbing/filter/cache"
)

const freqName = "freq."

func newFreqBuilder() *freqBuilder {
	return &freqBuilder{}
}

type freqBuilder struct{}

func (*freqBuilder) Name() string {
	return freqName
}

func (*freqBuilder) Build(name string) Variable {
	key := strings.TrimPrefix(name, freqName)
	if key == "" {
		return nil
	}

	return &FreqProfile{
		name: name,
		key:  key,
	}
}

// FreqProfile 频次控制
type FreqProfile struct {
	UnCacheableVariable
	name string
	key  string
}

func (s *FreqProfile) Name() string { return s.name }
func (s *FreqProfile) Value(ctx context.Context, data interface{}, _ *cache.Cache) (interface{}, error) {
	if getter, ok := data.(FrequencyGetter); ok {
		freData := getter.FrequencyValue(ctx, s.key)
		return freData, nil
	}
	return 0, nil
}
