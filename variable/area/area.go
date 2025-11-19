package area

import (
	"context"
	"errors"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/location"
	"github.com/liyanbing/filter/variable"
	"github.com/liyanbing/filter/variable/ip"
)

const (
	CountryName  = "country"
	ProvinceName = "province"
	CityName     = "city"
)

func init() {
	variable.Register(variable.NewSimpleVariable(&Area{
		name: CountryName,
	}))
	variable.Register(variable.NewSimpleVariable(&Area{
		name: ProvinceName,
	}))
	variable.Register(variable.NewSimpleVariable(&Area{
		name: CityName,
	}))
}

// Area 从IP中解析获取country信息
type Area struct {
	name string
}

func (s *Area) Cacheable() bool { return true }
func (s *Area) Name() string    { return s.name }
func (s *Area) Value(ctx context.Context, data interface{}, cache *cache.Cache) (interface{}, error) {
	ipVariable, ok := variable.Get(ip.Name)
	if !ok {
		return nil, errors.New("ip variable not found")
	}

	ipValue, err := variable.GetValue(ctx, ipVariable, data, cache)
	if err != nil {
		return nil, err
	}

	country, province, city, err := location.GetLocation(ipValue.(string))
	if err != nil {
		return nil, err
	}

	value := city
	switch s.name {
	case CountryName:
		value = country
	case ProvinceName:
		value = province
	}
	return value, nil
}
