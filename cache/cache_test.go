package cache

import (
	"os"
	"testing"
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
		if !ok {
			t.Logf("Get:%v not exists", v.Key)
			return
		}

		var ret interface{}
		ret = v.Value
		if ret != value {
			t.Logf("expected: %v but Got: %v", ret, value)
			return
		}
	}
}
