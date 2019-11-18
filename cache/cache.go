package cache

import "sync"

type Cache struct {
	Enable bool
	data   *sync.Map
	vars   *sync.Map
}

func (s *Cache) Set(key string, value interface{}) {
	if s == nil {
		return
	}

	if s.Enable {
		s.data.Store(key, value)
	}
}

func (s *Cache) Get(key string) (interface{}, bool) {
	if s == nil {
		return nil, false
	}

	elem, ok := s.data.Load(key)
	return elem, ok
}

func NewCache() *Cache {
	return &Cache{
		Enable: true,
		data:   &sync.Map{},
		vars:   &sync.Map{},
	}
}
