package variables

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Liyanbing/filter/cache"
	"github.com/Liyanbing/filter/location"
	"github.com/stretchr/testify/assert"

	filterContext "github.com/Liyanbing/filter/context"
)

func TestCalculator_GetName(t *testing.T) {
	cases := []struct {
		Reg string
		Ret bool
	}{
		{
			Reg: "get.name",
			Ret: true,
		},
		{
			Reg: "get.user{name}",
			Ret: true,
		},
		{
			Reg: "get.users{zhangsan}{name}",
			Ret: true,
		},
		{
			Reg: "get.users{zhangsan}[0]{name}",
			Ret: true,
		},
		{
			Reg: "get.user{name}[1]",
			Ret: true,
		},
		{
			Reg: "get.user[1]",
			Ret: true,
		},
		{
			Reg: "get.user[1]{name}",
			Ret: true,
		},
	}

	for _, v := range cases {
		filter := getReg.FindStringSubmatch(v.Reg)
		assert.Equal(t, 4, len(filter))
	}
}

const (
	_UserID    = "1"
	_Referer   = "http://www.baidu.com"
	_Channel   = "ios"
	_UserAgent = "Mozilla/5.0"
	_IP        = "47.107.69.99"
	_Platform  = "platform"
	_Device    = "device"
	_Version   = "version"
)

var (
	_Form = url.Values{
		"name": []string{"name"},
		"user": []string{`{"name":"name","age":18}`},
		"time": []string{"1"},
		"list": []string{"[1,2,3,4]"},
	}

	_UserTags = []string{"tag1", "tag2"}
)

func PrepayGeneralValues() filterContext.Values {
	return filterContext.Values{
		UserID:    _UserID,
		Referer:   _Referer,
		Channel:   _Channel,
		UserAgent: _UserAgent,
		IP:        _IP,
		GetForm:   _Form,
		Platform:  _Platform,
		Device:    _Device,
		Version:   _Version,
		UserTags:  _UserTags,
	}
}

func PrepayCustomData() map[string]interface{} {
	return map[string]interface{}{
		"name": "name",
		"age":  18,
	}
}

type CustomData struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (s *CustomData) CalcFactorGet(ctx context.Context, name string) (float64, error) {
	fmt.Println("获取计算值：", name)
	if name == "age" {
		return float64(s.Age), nil
	}

	return 1, nil
}

func (s *CustomData) CalcFactorSet(ctx context.Context, name string, value float64) {
	fmt.Println("计算设置：", name, value)
}

func (s *CustomData) FrequencyGet(ctx context.Context, name string) interface{} {
	fmt.Println("获取频率值：", name)
	return s.Age
}

func TestVariables(t *testing.T) {
	err := location.NewLocationWithDBFile("/Users/Leo/Desktop/GeoLite2-City/GeoLite2-City.mmdb")
	assert.Equal(t, nil, err)
	defer location.Close()

	ctx := context.Background()
	ctx = filterContext.WithContextValues(ctx, PrepayGeneralValues())
	ctx = filterContext.WithContextCustom(ctx, PrepayCustomData())
	ctx = filterContext.WithContextFilterID(ctx, "filter_id")

	now := time.Now()
	tsSimple, _ := strconv.ParseUint(now.Format("20060102150405"), 10, 64)
	customData := &CustomData{
		Name: "name",
		Age:  18,
	}

	cases := []struct {
		Variable string
		Name     string
		Assert   func(value interface{}) bool
	}{
		{
			Variable: "success",
			Name:     "success",
			Assert: func(value interface{}) bool {
				return value == 1
			},
		},
		{
			Variable: "rand",
			Name:     "rand",
			Assert: func(value interface{}) bool {
				v := value.(int)
				return v >= 1 && v <= 100
			},
		},
		{
			Variable: "ip",
			Name:     "ip",
			Assert: func(value interface{}) bool {
				return value == _IP
			},
		},
		{
			Variable: "country",
			Name:     "country",
			Assert: func(value interface{}) bool {
				return value == "中国"
			},
		},
		{
			Variable: "province",
			Name:     "province",
			Assert: func(value interface{}) bool {
				return value == "浙江省"
			},
		},
		{
			Variable: "city",
			Name:     "city",
			Assert: func(value interface{}) bool {
				return value == "杭州"
			},
		},
		{
			Variable: "timestamp",
			Name:     "timestamp",
			Assert: func(value interface{}) bool {
				return value == now.Unix()
			},
		},
		{
			Variable: "ts_simple",
			Name:     "ts_simple",
			Assert: func(value interface{}) bool {
				return value == tsSimple
			},
		},
		{
			Variable: "second",
			Name:     "second",
			Assert: func(value interface{}) bool {
				return value == now.Second()
			},
		},
		{
			Variable: "minute",
			Name:     "minute",
			Assert: func(value interface{}) bool {
				return value == now.Minute()
			},
		},
		{
			Variable: "hour",
			Name:     "hour",
			Assert: func(value interface{}) bool {
				return value == now.Hour()
			},
		},
		{
			Variable: "day",
			Name:     "day",
			Assert: func(value interface{}) bool {
				return value == now.Day()
			},
		},
		{
			Variable: "month",
			Name:     "month",
			Assert: func(value interface{}) bool {
				return value == int(now.Month())
			},
		},
		{
			Variable: "year",
			Name:     "year",
			Assert: func(value interface{}) bool {
				return value == now.Year()
			},
		},
		{
			Variable: "wday",
			Name:     "wday",
			Assert: func(value interface{}) bool {
				return value == int(now.Weekday())
			},
		},
		{
			Variable: "date",
			Name:     "date",
			Assert: func(value interface{}) bool {
				return value == now.Format("2006-01-02")
			},
		},
		{
			Variable: "time",
			Name:     "time",
			Assert: func(value interface{}) bool {
				return value == now.Format("2006-01-02 15:04:05")
			},
		},
		{
			Variable: "ua",
			Name:     "ua",
			Assert: func(value interface{}) bool {
				return value == _UserAgent
			},
		},
		{
			Variable: "referer",
			Name:     "referer",
			Assert: func(value interface{}) bool {
				return value == _Referer
			},
		},
		{
			Variable: "is_login",
			Name:     "is_login",
			Assert: func(value interface{}) bool {
				return value == true
			},
		},
		{
			Variable: "version",
			Name:     "version",
			Assert: func(value interface{}) bool {
				return value == _Version
			},
		},
		{
			Variable: "platform",
			Name:     "platform",
			Assert: func(value interface{}) bool {
				return value == _Platform
			},
		},
		{
			Variable: "channel",
			Name:     "channel",
			Assert: func(value interface{}) bool {
				return value == _Channel
			},
		},
		{
			Variable: "uid",
			Name:     "uid",
			Assert: func(value interface{}) bool {
				return value == _UserID
			},
		},
		{
			Variable: "device",
			Name:     "device",
			Assert: func(value interface{}) bool {
				return value == _Device
			},
		},
		{
			Variable: "user_tag",
			Name:     "user_tag",
			Assert: func(value interface{}) bool {
				ut1 := strings.Join(_UserTags, ",")
				ut2 := strings.Join(value.([]string), ",")
				return ut1 == ut2
			},
		},
		{
			Variable: "get.user{name}",
			Name:     "get.user{name}",
			Assert: func(value interface{}) bool {
				return value == "name"
			},
		},
		{
			Variable: "get.user{age}",
			Name:     "get.user{age}",
			Assert: func(value interface{}) bool {
				return value == float64(18)
			},
		},
		{
			Variable: "get.list[0]",
			Name:     "get.list[0]",
			Assert: func(value interface{}) bool {
				return value == float64(1)
			},
		},
		{
			Variable: "get.list[1]",
			Name:     "get.list[1]",
			Assert: func(value interface{}) bool {
				return value == float64(2)
			},
		},
		{
			Variable: "data.Name",
			Name:     "data.Name",
			Assert: func(value interface{}) bool {
				return value == "name"
			},
		},
		{
			Variable: "data.Age",
			Name:     "data.Age",
			Assert: func(value interface{}) bool {
				return value == 18

			},
		},
		{
			Variable: "calc.__age*__age",
			Name:     "calc.__age*__age",
			Assert: func(value interface{}) bool {
				return value == float64(18*18)
			},
		},
		{
			Variable: "calc.__age - __age",
			Name:     "calc.__age - __age",
			Assert: func(value interface{}) bool {
				return value == float64(18-18)
			},
		},
		{
			Variable: "freq.uid.daily",
			Name:     "freq.uid.daily",
			Assert: func(value interface{}) bool {
				return value == 18
			},
		},
		{
			Variable: "ctx.name",
			Name:     "ctx.name",
			Assert: func(value interface{}) bool {
				return value == "name"
			},
		},
		{
			Variable: "ctx.age",
			Name:     "ctx.age",
			Assert: func(value interface{}) bool {
				return value == 18
			},
		},
	}

	cacheService := cache.NewCache()
	for _, v := range cases {
		variable := Factory.Get(v.Variable)
		name := variable.GetName()
		assert.Equal(t, v.Name, name, v.Name)
		assert.Equal(t, true, v.Assert(variable.Value(ctx, customData, cacheService)), v.Name)
	}
}
