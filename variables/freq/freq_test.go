package freq

import (
	"context"
	"fmt"
	"testing"

	"github.com/airunny/filter/cache"
	_ "github.com/airunny/filter/location"
	"github.com/airunny/filter/variables"
	"github.com/stretchr/testify/assert"
)

type ValuerData struct {
	value interface{}
	err   error
}

func (s *ValuerData) FrequencyValue(ctx context.Context, key string) (interface{}, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.value, nil
}

func TestData(t *testing.T) {
	ctx := context.Background()
	cc := cache.NewCache()
	cases := []struct {
		Data interface{}
		Key  string
		Want interface{}
		Err  error
	}{
		// Valuer
		{
			Data: &ValuerData{
				value: "value",
			},
			Key:  "freq.name",
			Want: "value",
		},
		{
			Data: &ValuerData{
				value: "value",
				err:   fmt.Errorf("not found"),
			},
			Key:  "freq.name",
			Want: nil,
			Err:  fmt.Errorf("not found"),
		},
	}

	for index, tt := range cases {
		variable, ok := variables.Get(tt.Key)
		assert.True(t, ok)
		assert.NotNil(t, variable)
		assert.Equal(t, tt.Key, variable.Name(), index)

		ret, err := variable.Value(ctx, tt.Data, cc)
		if err != nil {
			assert.Equal(t, tt.Err, err)
		} else {
			assert.Equal(t, tt.Want, ret, index)
		}
	}
}
