package area

import (
	"context"
	"errors"

	"github.com/airunny/filter/cache"
	"github.com/airunny/filter/location"
	"github.com/airunny/filter/variables"
	"github.com/airunny/filter/variables/ip"
)

const (
	CountryName  = "country"
	ProvinceName = "province"
	CityName     = "city"
)

func init() {
	variables.Register(variables.NewSimpleVariable(&Area{
		name: CountryName,
	}))
	variables.Register(variables.NewSimpleVariable(&Area{
		name: ProvinceName,
	}))
	variables.Register(variables.NewSimpleVariable(&Area{
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
	ipVariable, ok := variables.Get(ip.Name)
	if !ok {
		return nil, errors.New("ip variable not found")
	}

	ipValue, err := variables.GetValue(ctx, ipVariable, data, cache)
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
