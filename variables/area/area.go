package area

import (
	"context"
	"errors"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/location"
	"github.com/liyanbing/filter/variables"
	"github.com/liyanbing/filter/variables/ip"
)

const (
	countryName  = "country"
	provinceName = "province"
	cityName     = "city"
)

func init() {
	variables.Register(variables.NewSimpleVariable(countryVariable()))
	variables.Register(variables.NewSimpleVariable(provinceVariable()))
	variables.Register(variables.NewSimpleVariable(cityVariable()))
}

func countryVariable() *Area {
	return &Area{
		name: countryName,
	}
}

func provinceVariable() *Area {
	return &Area{
		name: provinceName,
	}
}

func cityVariable() *Area {
	return &Area{
		name: cityName,
	}
}

// Area 从IP中解析获取country信息
type Area struct {
	variables.CacheableVariable
	name string
}

func (s *Area) Name() string { return s.name }

func (s *Area) Value(ctx context.Context, data interface{}, cache *cache.Cache) (interface{}, error) {
	variable, ok := variables.Get(ip.IPName)
	if !ok {
		return nil, errors.New("ip variable not found")
	}

	ip, err := variables.GetVariableValue(ctx, variable, data, cache)
	if err != nil {
		return nil, err
	}

	country, province, city, err := location.GetLocation(ip.(string))
	if err != nil {
		return nil, err
	}

	value := city
	switch s.name {
	case countryName:
		value = country
	case provinceName:
		value = province
	}
	return value, nil
}
