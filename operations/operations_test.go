package operations

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/location"
	"github.com/liyanbing/filter/variables"
	"github.com/stretchr/testify/assert"

	filterContext "github.com/liyanbing/filter/filter_context"
)

const (
	_UserID    = "1"
	_Referer   = "http://www.baidu.com"
	_Channel   = "ios"
	_UserAgent = "Mozilla/5.0"
	_IP        = "47.107.69.99"
	_Platform  = "platform"
	_Device    = "device"
	_Version   = "version"
	_Country   = "中国"
	_Province  = "浙江省"
	_City      = "杭州"
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

func PrepayGeneralValues() filterContext.CommonValue {
	return filterContext.CommonValue{
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

var (
	ctx        context.Context
	customData *CustomData
)

func TestMain(m *testing.M) {
	err := location.NewLocationWithDBFile("/Users/Leo/Desktop/GeoLite2-City/GeoLite2-City.mmdb")
	assert.Equal(&testing.T{}, nil, err)
	defer location.Close()

	ctx = context.Background()
	ctx = filterContext.WithCommonValue(ctx, PrepayGeneralValues())
	ctx = filterContext.WithCustom(ctx, PrepayCustomData())

	customData = &CustomData{
		Name: "name",
		Age:  18,
	}

	os.Exit(m.Run())
}

func TestEqual(t *testing.T) {
	now := time.Now()
	tsSimple, _ := strconv.ParseUint(now.Format("20060102150405"), 10, 64)

	cases := []struct {
		VariableName string
		Value        interface{}
		Ret          bool
	}{
		// success
		{
			VariableName: "success",
			Value:        1,
			Ret:          true,
		},
		{
			VariableName: "success",
			Value:        2,
			Ret:          false,
		},
		{
			VariableName: "rand",
			Value:        1,
			Ret:          false,
		},
		// ip
		{
			VariableName: "ip",
			Value:        _IP,
			Ret:          true,
		},
		{
			VariableName: "ip",
			Value:        _IP + "1",
			Ret:          false,
		},
		// country
		{
			VariableName: "country",
			Value:        _Country,
			Ret:          true,
		},
		{
			VariableName: "country",
			Value:        _Country + "1",
			Ret:          false,
		},
		// province
		{
			VariableName: "province",
			Value:        _Province,
			Ret:          true,
		},
		{
			VariableName: "province",
			Value:        _Province + "1",
			Ret:          false,
		},
		// city
		{
			VariableName: "city",
			Value:        _City,
			Ret:          true,
		},
		{
			VariableName: "city",
			Value:        _City + "1",
			Ret:          false,
		},
		// timestamp
		{
			VariableName: "timestamp",
			Value:        now.Unix(),
			Ret:          true,
		},
		{
			VariableName: "timestamp",
			Value:        now.Unix() + 10,
			Ret:          false,
		},
		// ts_simple
		{
			VariableName: "ts_simple",
			Value:        tsSimple,
			Ret:          true,
		},
		{
			VariableName: "ts_simple",
			Value:        tsSimple + 1,
			Ret:          false,
		},
		// second
		{
			VariableName: "second",
			Value:        now.Second(),
			Ret:          true,
		},
		{
			VariableName: "second",
			Value:        now.Second() + 1,
			Ret:          false,
		},
		// minute
		{
			VariableName: "minute",
			Value:        now.Minute(),
			Ret:          true,
		},
		{
			VariableName: "minute",
			Value:        now.Minute() + 1,
			Ret:          false,
		},
		// hour
		{
			VariableName: "hour",
			Value:        now.Hour(),
			Ret:          true,
		},
		{
			VariableName: "hour",
			Value:        now.Hour() + 1,
			Ret:          false,
		},
		// day
		{
			VariableName: "day",
			Value:        now.Day(),
			Ret:          true,
		},
		{
			VariableName: "day",
			Value:        now.Day() + 1,
			Ret:          false,
		},
		// month
		{
			VariableName: "month",
			Value:        int(now.Month()),
			Ret:          true,
		},
		{
			VariableName: "month",
			Value:        int(now.Month()) + 1,
			Ret:          false,
		},
		// year
		{
			VariableName: "year",
			Value:        now.Year(),
			Ret:          true,
		},
		{
			VariableName: "year",
			Value:        now.Year() + 1,
			Ret:          false,
		},
		// wday
		{
			VariableName: "wday",
			Value:        int(now.Weekday()),
			Ret:          true,
		},
		{
			VariableName: "wday",
			Value:        int(now.Weekday()) + 1,
			Ret:          false,
		},
		// date
		{
			VariableName: "date",
			Value:        now.Format("2006-01-02"),
			Ret:          true,
		},
		{
			VariableName: "date",
			Value:        now.Format("2006-01-02") + "1",
			Ret:          false,
		},
		// time
		{
			VariableName: "time",
			Value:        now.Format("2006-01-02 15:04:05"),
			Ret:          true,
		},
		{
			VariableName: "time",
			Value:        now.Format("2006-01-02 15:04:05") + "1",
			Ret:          false,
		},
		// ua
		{
			VariableName: "ua",
			Value:        _UserAgent,
			Ret:          true,
		},
		{
			VariableName: "ua",
			Value:        _UserAgent + "1",
			Ret:          false,
		},
		// referer
		{
			VariableName: "referer",
			Value:        _Referer,
			Ret:          true,
		},
		{
			VariableName: "referer",
			Value:        _Referer + "1",
			Ret:          false,
		},
		// is_login
		{
			VariableName: "is_login",
			Value:        true,
			Ret:          true,
		},
		{
			VariableName: "is_login",
			Value:        false,
			Ret:          false,
		},
		// version
		{
			VariableName: "version",
			Value:        _Version,
			Ret:          true,
		},
		{
			VariableName: "version",
			Value:        _Version + "1",
			Ret:          false,
		},
		// platform
		{
			VariableName: "platform",
			Value:        _Platform,
			Ret:          true,
		},
		{
			VariableName: "platform",
			Value:        _Platform + "1",
			Ret:          false,
		},
		// channel
		{
			VariableName: "channel",
			Value:        _Channel,
			Ret:          true,
		},
		{
			VariableName: "channel",
			Value:        _Channel + "1",
			Ret:          false,
		},
		// uid
		{
			VariableName: "uid",
			Value:        _UserID,
			Ret:          true,
		},
		{
			VariableName: "uid",
			Value:        _UserID + "1",
			Ret:          false,
		},
		// device
		{
			VariableName: "device",
			Value:        _Device,
			Ret:          true,
		},
		{
			VariableName: "device",
			Value:        _Device + "1",
			Ret:          false,
		},
		// user_tag
		{
			VariableName: "user_tag",
			Value:        _UserTags,
			Ret:          true,
		},
		{
			VariableName: "user_tag",
			Value:        []interface{}{"1", "2"},
			Ret:          false,
		},
		//get.user{name}
		{
			VariableName: "get.user{name}",
			Value:        "name",
			Ret:          true,
		},
		{
			VariableName: "get.user{name}",
			Value:        "name1",
			Ret:          false,
		},
		// get.user{age}
		{
			VariableName: "get.user{age}",
			Value:        18,
			Ret:          true,
		},
		{
			VariableName: "get.user{age}",
			Value:        19,
			Ret:          false,
		},
		// data.Name
		{
			VariableName: "data.Name",
			Value:        "name",
			Ret:          true,
		},
		{
			VariableName: "data.Name",
			Value:        "name1",
			Ret:          false,
		},
		// data.Age
		{
			VariableName: "data.Age",
			Value:        18,
			Ret:          true,
		},
		{
			VariableName: "data.Age",
			Value:        19,
			Ret:          false,
		},
		// calc.__age*__age
		{
			VariableName: "calc.__age*__age",
			Value:        18 * 18,
			Ret:          true,
		},
		{
			VariableName: "calc.__age*__age",
			Value:        19 * 19,
			Ret:          false,
		},
		// freq.uid.daily
		{
			VariableName: "freq.uid.daily",
			Value:        18,
			Ret:          true,
		},
		{
			VariableName: "freq.uid.daily",
			Value:        19,
			Ret:          false,
		},
		// ctx.name
		{
			VariableName: "ctx.name",
			Value:        "name",
			Ret:          true,
		},
		{
			VariableName: "ctx.name",
			Value:        "name1",
			Ret:          false,
		},
		// ctx.age
		{
			VariableName: "ctx.age",
			Value:        18,
			Ret:          true,
		},
		{
			VariableName: "ctx.age",
			Value:        19,
			Ret:          false,
		},
	}

	cacheService := cache.NewCache()
	for _, v := range cases {
		variable := variables.Factory.Get(v.VariableName)
		operation := Factory.Get("=")
		prepayValue, err := operation.PrepareValue(v.Value)
		assert.Equal(t, nil, err, v.VariableName)
		assert.Equal(t, v.Ret, operation.Run(ctx, variable, prepayValue, customData, cacheService), v.VariableName)
	}
}
