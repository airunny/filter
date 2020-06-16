package cache

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var s *Cache

func TestMain(m *testing.M) {
	s = NewCache()
	os.Exit(m.Run())
}

func TestCache_Set(t *testing.T) {
	type Temp struct {
	}

	cases := []struct {
		Key   string
		Value interface{}
	}{
		{
			Key:   "1",
			Value: 1,
		},
		{
			Key:   "2",
			Value: "1",
		},
		{
			Key:   "3",
			Value: 1.01,
		},
		{
			Key:   "4",
			Value: true,
		},
		{
			Key:   "5",
			Value: Temp{},
		},
		{
			Key:   "6",
			Value: &Temp{},
		},
	}

	for _, v := range cases {
		s.Set(v.Key, v.Value)
		value, ok := s.Get(v.Key)
		assert.Equal(t, true, ok)
		assert.Equal(t, value, v.Value)
	}
}
