package variables

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/liyanbing/filter/cache"
	filterContext "github.com/liyanbing/filter/context"
	"github.com/stretchr/testify/assert"
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
			Reg: "get.user.name",
			Ret: true,
		},
		{
			Reg: "get.users.0.name",
			Ret: true,
		},
		{
			Reg: "get.users.zhangsan.0.name",
			Ret: true,
		},
		{
			Reg: "get.user.name.1",
			Ret: true,
		},
		{
			Reg: "get.user.1",
			Ret: true,
		},
		{
			Reg: "get.user.1.name",
			Ret: true,
		},
		{
			Reg: "get.user.1.name.1",
			Ret: true,
		},
		{
			Reg: "get.user.1.name.1.1.name.age",
			Ret: true,
		},
	}

	for _, v := range cases {
		filter := getReg.FindStringSubmatch(v.Reg)
		t.Log(filter)
		assert.Equal(t, 2, len(filter))
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

type CustomData struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (s *CustomData) CalcValue(ctx context.Context, name string) (float64, error) {
	if name == "age" {
		return float64(s.Age), nil
	}

	if name == "second" {
		return float64(18), nil
	}

	return 1, nil
}

func (s *CustomData) CalcFactorSet(ctx context.Context, name string, value float64) {
	fmt.Println("计算设置：", name, value)
}

func (s *CustomData) FrequencyValue(ctx context.Context, name string) interface{} {
	fmt.Println("获取频率值：", name)
	return s.Age
}

func (s *CustomData) Value(ctx context.Context, key string) interface{} {
	fmt.Println("获取值：", key)
	switch key {
	case "Name":
		return s.Name
	case "Age":
		return s.Age
	default:
		return nil
	}
}

func TestVariables(t *testing.T) {
	//err := location.NewLocationWithDBFile("/Users/Leo/Desktop/GeoLite2-City/GeoLite2-City.mmdb")
	//assert.Equal(t, nil, err)
	//defer location.Close()

	ctx := context.Background()
	ctx = filterContext.WithUserID(ctx, _UserID)
	ctx = filterContext.WithDevice(ctx, _Device)
	ctx = filterContext.WithIP(ctx, _IP)
	ctx = filterContext.WithVersion(ctx, _Version)
	ctx = filterContext.WithPlatform(ctx, _Platform)
	ctx = filterContext.WithChannel(ctx, _Channel)
	ctx = filterContext.WithUA(ctx, _UserAgent)
	ctx = filterContext.WithReferer(ctx, _Referer)
	ctx = filterContext.WithUserTag(ctx, _UserTags)
	ctx = context.WithValue(ctx, "name", "name")
	ctx = context.WithValue(ctx, "age", 18)

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
		//{
		//	Variable: "country",
		//	Name:     "country",
		//	Assert: func(value interface{}) bool {
		//		return value == "中国"
		//	},
		//},
		//{
		//	Variable: "province",
		//	Name:     "province",
		//	Assert: func(value interface{}) bool {
		//		return value == "浙江省"
		//	},
		//},
		//{
		//	Variable: "city",
		//	Name:     "city",
		//	Assert: func(value interface{}) bool {
		//		return value == "杭州"
		//	},
		//},
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
			Variable: "calc.__age * __second",
			Name:     "calc.__age * __second",
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
		variable, ok := Get(v.Variable)
		if !ok {
			continue
		}
		name := variable.Name()
		assert.Equal(t, v.Name, name, v.Name)
		value, err := variable.Value(ctx, customData, cacheService)
		assert.Nil(t, err, v.Name)
		assert.Equal(t, true, v.Assert(value), v.Name)
	}
}
