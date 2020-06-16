package variables

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/liyanbing/calc/compute"
	"github.com/liyanbing/calc/variables"
	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/location"
	"github.com/liyanbing/filter/utils"

	filterContext "github.com/liyanbing/filter/filter_context"
	filterType "github.com/liyanbing/filter/filter_type"
)

var getReg = regexp.MustCompile(`^get.(.+)`)

type Variable interface {
	Cacheable() bool
	Value(ctx context.Context, data interface{}, cache *cache.Cache) interface{}
	GetName() string
}

type CacheableVariable struct{}

func (s *CacheableVariable) Cacheable() bool {
	return true
}

type UnCacheableVariable struct{}

func (s *UnCacheableVariable) Cacheable() bool {
	return false
}

type Creator func(string) Variable

func simpleVariableCreator(instance Variable) Creator {
	return func(name string) Variable {
		return instance
	}
}

// ---------------factory--------------
type factory struct {
	creators map[string]Creator
}

func (s *factory) Get(name string) Variable {
	if creator, ok := s.creators[name]; ok {
		return creator(name)
	} else {
		segments := strings.Split(name, ".")
		if len(segments) > 1 {
			if creator, ok := s.creators[segments[0]+"."]; ok {
				return creator(name)
			}
		}
	}

	return nil
}

func (s *factory) Register(name string, creator Creator) error {
	if _, ok := s.creators[name]; ok {
		return fmt.Errorf("%v variable already exists", name)
	}
	s.creators[name] = creator
	return nil
}

func RegisterVariable(name string, creator Creator) error {
	return Factory.Register(name, creator)
}

var _ = RegisterVariable

func RegisterVariableFunc(name string, vf FuncNoCache) error {
	return Factory.Register(name, simpleVariableCreator(vf))
}

type FuncNoCache func(ctx context.Context, data interface{}, cache *cache.Cache) interface{}

func (vf FuncNoCache) Cacheable() bool {
	return false
}

func (vf FuncNoCache) GetName() string {
	return ""
}

func (vf FuncNoCache) Value(ctx context.Context, data interface{}, cache *cache.Cache) interface{} {
	return vf(ctx, data, cache)
}

// ---------------general variables register--------------
var Factory = &factory{
	creators: make(map[string]Creator),
}

func init() {
	_ = Factory.Register("success", simpleVariableCreator(&Success{}))
	_ = Factory.Register("rand", simpleVariableCreator(&Rand{}))
	_ = Factory.Register("ip", simpleVariableCreator(&IP{}))
	_ = Factory.Register("country", simpleVariableCreator(&Area{name: "country"}))
	_ = Factory.Register("province", simpleVariableCreator(&Area{name: "province"}))
	_ = Factory.Register("city", simpleVariableCreator(&Area{name: "city"}))
	_ = Factory.Register("timestamp", simpleVariableCreator(&Time{name: "timestamp"}))
	_ = Factory.Register("ts_simple", simpleVariableCreator(&Time{name: "ts_simple"}))
	_ = Factory.Register("second", simpleVariableCreator(&Time{name: "second"}))
	_ = Factory.Register("minute", simpleVariableCreator(&Time{name: "minute"}))
	_ = Factory.Register("hour", simpleVariableCreator(&Time{name: "hour"}))
	_ = Factory.Register("day", simpleVariableCreator(&Time{name: "day"}))
	_ = Factory.Register("month", simpleVariableCreator(&Time{name: "month"}))
	_ = Factory.Register("year", simpleVariableCreator(&Time{name: "year"}))
	_ = Factory.Register("wday", simpleVariableCreator(&Time{name: "wday"}))
	_ = Factory.Register("date", simpleVariableCreator(&Time{name: "date"}))
	_ = Factory.Register("time", simpleVariableCreator(&Time{name: "time"}))
	_ = Factory.Register("ua", simpleVariableCreator(&UserAgent{}))
	_ = Factory.Register("referer", simpleVariableCreator(&Referer{}))
	_ = Factory.Register("is_login", simpleVariableCreator(&IsLogin{}))
	_ = Factory.Register("version", simpleVariableCreator(&Version{}))
	_ = Factory.Register("platform", simpleVariableCreator(&Platform{}))
	_ = Factory.Register("channel", simpleVariableCreator(&Channel{}))
	_ = Factory.Register("uid", simpleVariableCreator(&UID{}))
	_ = Factory.Register("device", simpleVariableCreator(&Device{}))
	_ = Factory.Register("user_tag", simpleVariableCreator(&UserTag{}))
	_ = Factory.Register("get.", GetCreator)
	_ = Factory.Register("data.", DataCreator)
	_ = Factory.Register("calc.", CalculatorCreator)
	_ = Factory.Register("freq.", FreqProfileCreator)
	_ = Factory.Register("ctx.", CtxCreator)
}

func GetVariableValue(ctx context.Context, v Variable, data interface{}, cache *cache.Cache) interface{} {
	if v == nil {
		return ""
	}

	if v.Cacheable() {
		if value, ok := cache.Get(v.GetName()); ok {
			return value
		}
	}

	value := v.Value(ctx, data, cache)
	if v.Cacheable() {
		cache.Set(v.GetName(), value)
	}

	return value
}

// ----------------variables-------------
// success
type Success struct{ CacheableVariable }

func (s *Success) GetName() string { return "success" }
func (s *Success) Value(_ context.Context, _ interface{}, _ *cache.Cache) interface{} {
	return 1
}

// rand
type Rand struct{ UnCacheableVariable }

func (s *Rand) GetName() string { return "rand" }
func (s *Rand) Value(_ context.Context, _ interface{}, _ *cache.Cache) interface{} {
	return rand.Intn(100) + 1
}

// ip
type IP struct{ CacheableVariable }

func (s *IP) GetName() string { return "ip" }
func (s *IP) Value(ctx context.Context, _ interface{}, _ *cache.Cache) interface{} {
	ip, _ := filterContext.IP(ctx)
	return ip
}

// area
type Area struct {
	CacheableVariable
	name string
}

func (s *Area) GetName() string { return s.name }
func (s *Area) Value(ctx context.Context, data interface{}, cache *cache.Cache) interface{} {
	ip := GetVariableValue(ctx, Factory.Get("ip"), data, cache).(string)
	country, province, city, err := location.GetLocation(ip)
	if err != nil {
		panic(err)
	}

	dataValue := city
	switch s.name {
	case "country":
		dataValue = country
	case "province":
		dataValue = province
	}
	return dataValue
}

// time
type Time struct {
	CacheableVariable
	name string
}

func (s *Time) GetName() string { return s.name }
func (s *Time) Value(_ context.Context, _ interface{}, _ *cache.Cache) interface{} {
	now := time.Now()

	switch s.name {
	case "timestamp":
		return now.Unix()
	case "ts_simple":
		ret, _ := strconv.ParseUint(now.Format("20060102150405"), 10, 64)
		return ret
	case "second":
		return now.Second()
	case "minute":
		return now.Minute()
	case "hour":
		return now.Hour()
	case "day":
		return now.Day()
	case "month":
		return int(now.Month())
	case "year":
		return now.Year()
	case "wday":
		return int(now.Weekday())
	case "date":
		return now.Format("2006-01-02")
	default:
		return now.Format("2006-01-02 15:04:05") // time
	}
}

// user agent
type UserAgent struct{ CacheableVariable }

func (s *UserAgent) GetName() string { return "ua" }
func (s *UserAgent) Value(ctx context.Context, _ interface{}, _ *cache.Cache) interface{} {
	ua, _ := filterContext.UserAgent(ctx)
	return ua
}

// referer
type Referer struct{ CacheableVariable }

func (s *Referer) GetName() string { return "referer" }
func (s *Referer) Value(ctx context.Context, _ interface{}, _ *cache.Cache) interface{} {
	url, _ := filterContext.Referer(ctx)
	return url
}

// is_login
type IsLogin struct{ CacheableVariable }

func (s *IsLogin) GetName() string { return "is_login" }
func (s *IsLogin) Value(ctx context.Context, _ interface{}, _ *cache.Cache) interface{} {
	uid, _ := filterContext.UserID(ctx)
	if uid != "" {
		return true
	} else {
		return false
	}
}

// version
type Version struct{ CacheableVariable }

func (s *Version) GetName() string { return "version" }
func (s *Version) Value(ctx context.Context, _ interface{}, _ *cache.Cache) interface{} {
	version, _ := filterContext.Version(ctx)
	return version
}

// platform
type Platform struct{ CacheableVariable }

func (s *Platform) GetName() string { return "platform" }
func (s *Platform) Value(ctx context.Context, _ interface{}, _ *cache.Cache) interface{} {
	plt, _ := filterContext.Platform(ctx)
	return plt
}

// channel
type Channel struct{ CacheableVariable }

func (s *Channel) GetName() string { return "channel" }
func (s *Channel) Value(ctx context.Context, _ interface{}, _ *cache.Cache) interface{} {
	chl, _ := filterContext.Channel(ctx)
	return chl
}

// uid
type UID struct{ CacheableVariable }

func (s *UID) GetName() string { return "uid" }
func (s *UID) Value(ctx context.Context, _ interface{}, _ *cache.Cache) interface{} {
	uid, _ := filterContext.UserID(ctx)
	return uid
}

// device
type Device struct{ CacheableVariable }

func (s *Device) GetName() string { return "device" }
func (s *Device) Value(ctx context.Context, _ interface{}, _ *cache.Cache) interface{} {
	device, _ := filterContext.Device(ctx)
	return device
}

// user_tag
type UserTag struct{ CacheableVariable }

func (s *UserTag) GetName() string { return "user_tag" }
func (s *UserTag) Value(ctx context.Context, _ interface{}, _ *cache.Cache) interface{} {
	tags, _ := filterContext.UserTags(ctx)
	return tags
}

//get
type Get struct {
	CacheableVariable
	name      string
	paramName string
}

func (s *Get) GetName() string { return s.name }
func (s *Get) Value(ctx context.Context, _ interface{}, _ *cache.Cache) interface{} {
	var formValue string
	values, ok := filterContext.Form(ctx)

	paramFilter := strings.Split(s.paramName, ".")
	if ok && len(paramFilter) > 0 {
		formValue = values.Get(paramFilter[0])
	}

	if formValue == "" {
		return formValue
	}

	var valueData interface{}
	if err := json.Unmarshal([]byte(formValue), &valueData); err != nil {
		return ""
	}

	if v, ok := utils.GetObjectValueByKey(valueData, strings.TrimPrefix(s.paramName, paramFilter[0]+".")); ok {
		return v
	} else {
		return ""
	}

}

func GetCreator(name string) Variable {
	if ma := getReg.FindStringSubmatch(name); len(ma) == 2 {
		return &Get{
			name:      ma[0],
			paramName: ma[1],
		}
	}

	return nil
}

// data
type Data struct {
	UnCacheableVariable
	name string
	key  string
}

func (s *Data) GetName() string { return s.name }
func (s *Data) Value(_ context.Context, data interface{}, _ *cache.Cache) interface{} {
	if v, ok := utils.GetObjectValueByKey(data, s.key); ok {
		return v
	}

	return ""
}

func DataCreator(name string) Variable {
	key := strings.TrimPrefix(name, "data.")

	if key == "" {
		return nil
	}

	return &Data{
		name: name,
		key:  key,
	}
}

// calculator
type Calculator struct {
	UnCacheableVariable
	name string
	expr string
}

func (s *Calculator) GetName() string { return s.name }
func (s *Calculator) Value(ctx context.Context, data interface{}, cache *cache.Cache) interface{} {
	value, _ := compute.Evaluate(s.expr, variables.ValueSourceFunc(func(name string) float64 {
		if getter, ok := data.(CalcFactorGetter); ok {
			v, err := getter.CalcFactorGet(ctx, name)
			if err == nil {
				return v
			}
		}

		variable := Factory.Get(name)
		if variable != nil {
			return filterType.GetFloat(variable.Value(ctx, data, cache))
		}

		return 0
	}))
	return value
}

func CalculatorCreator(name string) Variable {
	expr := strings.TrimPrefix(name, "calc.")

	if expr == "" {
		return nil
	}

	return &Calculator{
		name: name,
		expr: expr,
	}
}

// freq profile
type FreqProfile struct {
	UnCacheableVariable
	name string
	key  string
}

func (s *FreqProfile) GetName() string { return s.name }
func (s *FreqProfile) Value(ctx context.Context, data interface{}, _ *cache.Cache) interface{} {
	if getter, ok := data.(FrequencyGetter); ok {
		freData := getter.FrequencyGet(ctx, s.key)

		return freData
	}

	return ""
}

func FreqProfileCreator(name string) Variable {
	key := strings.TrimPrefix(name, "freq.")

	if key == "" {
		return nil
	}

	return &FreqProfile{
		name: name,
		key:  key,
	}
}

// ctx of custom
type Ctx struct {
	UnCacheableVariable
	name string
	key  string
}

func (s *Ctx) GetName() string { return s.name }
func (s *Ctx) Value(ctx context.Context, _ interface{}, _ *cache.Cache) interface{} {
	variableData, ok := filterContext.FromCustom(ctx)
	if !ok {
		return ""
	}

	if value, ok := variableData[s.key]; ok {
		return value
	}

	return ""
}

func CtxCreator(name string) Variable {
	key := strings.TrimPrefix(name, "ctx.")
	if key == "" {
		return nil
	}

	return &Ctx{
		name: name,
		key:  key,
	}
}

// ----------
type ITest interface {
	Run(ctx context.Context, id string) (string, error)
}

type TestArgsFunc func(ctx context.Context, id string) (string, error)

func (f TestArgsFunc) Run(ctx context.Context, id string) (string, error) {
	return f(ctx, id)
}

func Test(ctx context.Context, in ITest) (string, error) {
	return in.Run(ctx, "1")
}

func Func1(_ context.Context, id string) (string, error) {
	return id, nil
}
