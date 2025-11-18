package utils

import (
	"context"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
)

var EPSILON = 0.00000001

func FloatEquals(a, b float64) bool {
	if (a-b) < EPSILON && (b-a) < EPSILON {
		return true
	}

	return false
}

type IWeight interface {
	GetWeight() int64
}

func TotalWeight(weight []IWeight) int64 {
	total := int64(0)

	for _, v := range weight {
		total += v.GetWeight()
	}

	return total
}

func PickByWeight(weight []IWeight, totalWeight int64) int {
	if totalWeight == 0 {
		totalWeight = TotalWeight(weight)
	}

	choose := rand.Int63n(totalWeight) + 1
	line := int64(0)
	for i, b := range weight {
		line += b.GetWeight()
		if choose <= line {
			return i
		}
	}
	return 0
}

func ShuffleByWeight(weight []IWeight, totalWeight int64) {
	if len(weight) == 0 || len(weight) == 1 {
		return
	}

	if totalWeight == 0 {
		totalWeight = TotalWeight(weight)
	}

	for curIndex := 0; curIndex < len(weight); curIndex++ {
		chooseIndex := curIndex + PickByWeight(weight[curIndex:], totalWeight)
		weight[chooseIndex], weight[curIndex] = weight[curIndex], weight[chooseIndex]
		totalWeight -= weight[curIndex].GetWeight()
	}
}

func GetObjectValueByKey(ctx context.Context, data interface{}, key string) (interface{}, bool) {
	type Valuer interface {
		Value(ctx context.Context, key string) interface{}
	}
	if valuer, ok := data.(Valuer); ok {
		return valuer.Value(ctx, key), true
	}

	key = strings.TrimSpace(key)
	if key == "." || key == "" {
		return data, true
	}

	segs := strings.Split(key, ".")
	if len(segs) <= 0 {
		return data, true
	}

	for index := 0; index < len(segs); {
		seg := segs[index]

		if data == nil {
			return nil, false
		}
		seg = strings.TrimSpace(seg)

		switch reflect.TypeOf(data).Kind() {
		case reflect.Map:
			if v, ok := data.(map[string]interface{}); ok {
				if data, ok = v[seg]; !ok {
					return nil, false
				}
			} else {
				return nil, false
			}

		case reflect.Array, reflect.Slice:
			value := reflect.ValueOf(data)
			if i, err := strconv.Atoi(seg); err != nil {
				return nil, false
			} else if i < 0 || i >= value.Len() {
				return nil, false
			} else {
				data = value.Index(i).Interface()
			}

		case reflect.Struct:
			value := reflect.ValueOf(data)
			f := value.FieldByName(seg)
			if !f.IsValid() {
				return nil, false
			}
			data = f.Interface()
		case reflect.Ptr:
			data = reflect.ValueOf(data).Elem().Interface()
			continue

		default:
			return nil, false
		}
		index++
	}
	return data, true
}
