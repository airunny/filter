package variables

import (
	"context"
	"encoding/json"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Liyanbing/filter/cache"
	"github.com/Liyanbing/filter/location"
	"github.com/Liyanbing/filter/utils"
	"github.com/skOak/calc/compute"
	"github.com/skOak/calc/variables"

	filterContext "github.com/Liyanbing/filter/context"
	filterType "github.com/Liyanbing/filter/type"
)

var getReg = regexp.MustCompile(`^get.(.+?)(?:\{([^\}]+)\})?(?:\[(\d+)\])?$`)

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

// --------------data getter and setter---------------
type CalcFactorGetter interface {
	CalcFactorGet(ctx context.Context, name string) (float64, error)
}

type CalcFactorSetter interface {
	CalcFactorSet(ctx context.Context, name string, value float64)
}

type FrequencyGetter interface {
	FrequencyGet(ctx context.Context, name string) interface{}
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

func (s *factory) Register(name string, creator Creator) {
	s.creators[name] = creator
}

func RegisterVariable(name string, creator Creator) {
	Factory.Register(name, creator)
}

func RegisterVariableFunc(name string, vf FuncNoCache) {
	Factory.Register(name, simpleVariableCreator(vf))
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
	Factory.Register("success", simpleVariableCreator(&Success{}))
	Factory.Register("rand", simpleVariableCreator(&Rand{}))
	Factory.Register("ip", simpleVariableCreator(&IP{}))
	Factory.Register("country", simpleVariableCreator(&Area{name: "country"}))
	Factory.Register("province", simpleVariableCreator(&Area{name: "province"}))
	Factory.Register("city", simpleVariableCreator(&Area{name: "city"}))
	Factory.Register("timestamp", simpleVariableCreator(&Time{name: "timestamp"}))
	Factory.Register("ts_simple", simpleVariableCreator(&Time{name: "ts_simple"}))
	Factory.Register("second", simpleVariableCreator(&Time{name: "second"}))
	Factory.Register("minute", simpleVariableCreator(&Time{name: "minute"}))
	Factory.Register("hour", simpleVariableCreator(&Time{name: "hour"}))
	Factory.Register("day", simpleVariableCreator(&Time{name: "day"}))
	Factory.Register("month", simpleVariableCreator(&Time{name: "month"}))
	Factory.Register("year", simpleVariableCreator(&Time{name: "year"}))
	Factory.Register("wday", simpleVariableCreator(&Time{name: "wday"}))
	Factory.Register("date", simpleVariableCreator(&Time{name: "date"}))
	Factory.Register("time", simpleVariableCreator(&Time{name: "time"}))
	Factory.Register("ua", simpleVariableCreator(&UserAgent{}))
	Factory.Register("referer", simpleVariableCreator(&Referer{}))
	Factory.Register("is_login", simpleVariableCreator(&IsLogin{}))
	Factory.Register("version", simpleVariableCreator(&Version{}))
	Factory.Register("platform", simpleVariableCreator(&Platform{}))
	Factory.Register("channel", simpleVariableCreator(&Channel{}))
	Factory.Register("uid", simpleVariableCreator(&UID{}))
	Factory.Register("device", simpleVariableCreator(&Device{}))
	Factory.Register("user_tag", simpleVariableCreator(&UserTag{}))
	Factory.Register("get.", GetCreator)
	Factory.Register("data.", DataCreator)
	Factory.Register("calc.", CalculatorCreator)
	Factory.Register("freq.", FreqProfileCreator)
	Factory.Register("ctx.", CtxCreator)
}

func GetVariableValue(v Variable, ctx context.Context, data interface{}, cache *cache.Cache) interface{} {
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
func (s *Success) Value(ctx context.Context, data interface{}, cache *cache.Cache) interface{} {
	return 1
}

// rand
type Rand struct{ UnCacheableVariable }

func (s *Rand) GetName() string { return "rand" }
func (s *Rand) Value(ctx context.Context, data interface{}, cache *cache.Cache) interface{} {
	return rand.Intn(100) + 1
}

// ip
type IP struct{ CacheableVariable }

func (s *IP) GetName() string { return "ip" }
func (s *IP) Value(ctx context.Context, data interface{}, cache *cache.Cache) interface{} {
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
	ip := GetVariableValue(Factory.Get("ip"), ctx, data, cache).(string)
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
func (s *Time) Value(ctx context.Context, data interface{}, cache *cache.Cache) interface{} {
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
func (s *UserAgent) Value(ctx context.Context, data interface{}, cache *cache.Cache) interface{} {
	ua, _ := filterContext.UserAgent(ctx)
	return ua
}

// referer
type Referer struct{ CacheableVariable }

func (s *Referer) GetName() string { return "referer" }
func (s *Referer) Value(ctx context.Context, data interface{}, cache *cache.Cache) interface{} {
	url, _ := filterContext.Referer(ctx)
	return url
}

// is_login
type IsLogin struct{ CacheableVariable }

func (s *IsLogin) GetName() string { return "is_login" }
func (s *IsLogin) Value(ctx context.Context, data interface{}, cache *cache.Cache) interface{} {
	uid, _ := filterContext.UserAgent(ctx)
	if uid != "" {
		return true
	} else {
		return false
	}
}

// version
type Version struct{ CacheableVariable }

func (s *Version) GetName() string { return "version" }
func (s *Version) Value(ctx context.Context, data interface{}, cache *cache.Cache) interface{} {
	version, _ := filterContext.Version(ctx)
	return version
}

// platform
type Platform struct{ CacheableVariable }

func (s *Platform) GetName() string { return "platform" }
func (s *Platform) Value(ctx context.Context, data interface{}, cache *cache.Cache) interface{} {
	plt, _ := filterContext.Platform(ctx)
	return plt
}

// channel
type Channel struct{ CacheableVariable }

func (s *Channel) GetName() string { return "channel" }
func (s *Channel) Value(ctx context.Context, data interface{}, cache *cache.Cache) interface{} {
	chl, _ := filterContext.Channel(ctx)
	return chl
}

// uid
type UID struct{ CacheableVariable }

func (s *UID) GetName() string { return "uid" }
func (s *UID) Value(ctx context.Context, data interface{}, cache *cache.Cache) interface{} {
	uid, _ := filterContext.UserID(ctx)
	return uid
}

// device
type Device struct{ CacheableVariable }

func (s *Device) GetName() string { return "device" }
func (s *Device) Value(ctx context.Context, data interface{}, cache *cache.Cache) interface{} {
	device, _ := filterContext.Device(ctx)
	return device
}

// user_tag
type UserTag struct{ CacheableVariable }

func (s *UserTag) GetName() string { return "user_tag" }
func (s *UserTag) Value(ctx context.Context, data interface{}, cache *cache.Cache) interface{} {
	tags, _ := filterContext.UserTags(ctx)
	return tags
}

//get
type Get struct {
	CacheableVariable
	name      string
	paramName string
	listMode  bool
	listIndex int
	jsonMode  bool
	jsonKey   string
}

func (s *Get) GetName() string { return s.name }
func (s *Get) Value(ctx context.Context, data interface{}, cache *cache.Cache) interface{} {
	var formValue string
	values, ok := filterContext.Form(ctx)
	if ok {
		formValue = values.Get(s.paramName)
	}

	if formValue == "" || (!s.listMode && !s.jsonMode) {
		return formValue
	}

	var valueData interface{}
	if err := json.Unmarshal([]byte(formValue), &valueData); err != nil {
		return ""
	}

	if data == nil {
		return ""
	}

	var result interface{}
	if s.jsonMode {
		if v, ok := utils.GetObjectValueByKey(valueData, s.jsonKey); ok {
			result = v
		} else {
			return ""
		}
	}

	if s.listMode {
		tp := filterType.GetMyType(valueData)
		if tp == filterType.STRING {
			values := strings.Split(valueData.(string), ",")
			if s.listIndex < 0 || s.listIndex >= len(values) {
				return ""
			}
			result = values[s.listIndex]

		} else if tp == filterType.ARRAY {
			values := valueData.([]interface{})
			if s.listIndex < 0 || s.listIndex >= len(values) {
				return ""
			}
			result = values[s.listIndex]
		}
	}

	return result
}

func GetCreator(name string) Variable {
	if ma := getReg.FindStringSubmatch(name); len(ma) == 4 {
		obj := &Get{
			name:      ma[0],
			paramName: ma[1],
			listMode:  false,
			listIndex: 0,
			jsonMode:  false,
			jsonKey:   "",
		}

		if ma[2] != "" {
			obj.jsonMode = true
			obj.jsonKey = ma[2]
		}

		if ma[3] != "" {
			obj.listMode = true
			obj.listIndex, _ = strconv.Atoi(ma[3])
		}

		return obj
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
func (s *Data) Value(ctx context.Context, data interface{}, cache *cache.Cache) interface{} {
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
			// record name&value in data if possible
			v := filterType.GetFloat(variable.Value(ctx, data, cache))
			if setter, ok := data.(CalcFactorSetter); ok {
				setter.CalcFactorSet(ctx, name, v)
			}
			return v
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
func (s *FreqProfile) Value(ctx context.Context, data interface{}, cache *cache.Cache) interface{} {
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
func (s *Ctx) Value(ctx context.Context, data interface{}, cache *cache.Cache) interface{} {
	variableData, ok := filterContext.FromContextCustom(ctx)
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
