package data

import (
	"context"
	"fmt"
	"testing"

	"github.com/liyanbing/filter/cache"
	_ "github.com/liyanbing/filter/location"
	"github.com/liyanbing/filter/variables"
	"github.com/stretchr/testify/assert"
)

type Work struct {
	WorkName string
}

type User struct {
	Name   string
	Age    int8
	IDCard string
	Works  []*Work
}

type Temp struct {
	User  *User
	Tests []string
	Merge map[string]interface{}
}

type ValuerData struct {
	value interface{}
	err   error
}

func (s *ValuerData) Value(ctx context.Context, key string) (interface{}, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.value, nil
}

func TestData(t *testing.T) {
	ctx := context.Background()
	cc := cache.NewCache()
	cases := []struct {
		Data interface{}
		Key  string
		Want interface{}
		Err  error
	}{
		// Valuer
		{
			Data: &ValuerData{
				value: "value",
			},
			Key:  "data.name",
			Want: "value",
		},
		{
			Data: &ValuerData{
				value: "value",
				err:   fmt.Errorf("not found"),
			},
			Key:  "data.name",
			Want: nil,
			Err:  fmt.Errorf("not found"),
		},
		// map
		{
			Data: map[string]interface{}{
				"user": &Temp{
					User: &User{
						Name:   "zhangsan",
						Age:    18,
						IDCard: "110",
						Works: []*Work{
							{
								WorkName: "operation",
							},
						},
					},
					Tests: []string{"1", "2"},
				},
				"name":  "张三",
				"age":   18,
				"citys": []string{"1", "2"},
				"ages":  []int64{1, 2},
			},
			Key:  "data.name",
			Want: "张三",
		},
		{
			Data: map[string]interface{}{
				"user": &Temp{
					User: &User{
						Name:   "zhangsan",
						Age:    18,
						IDCard: "110",
						Works: []*Work{
							{
								WorkName: "operation",
							},
						},
					},
					Tests: []string{"1", "2"},
				},
				"name":  "张三",
				"age":   18,
				"citys": []string{"1", "2"},
				"ages":  []int64{1, 2},
			},
			Key:  "data.age",
			Want: 18,
		},
		{
			Data: map[string]interface{}{
				"user": &Temp{
					User: &User{
						Name:   "zhangsan",
						Age:    18,
						IDCard: "110",
						Works: []*Work{
							{
								WorkName: "operation",
							},
						},
					},
					Tests: []string{"1", "2"},
				},
				"name":  "张三",
				"age":   18,
				"citys": []string{"1", "2"},
				"ages":  []int64{1, 2},
			},
			Key:  "data.citys",
			Want: []string{"1", "2"},
		},
		{
			Data: map[string]interface{}{
				"user": &Temp{
					User: &User{
						Name:   "zhangsan",
						Age:    18,
						IDCard: "110",
						Works: []*Work{
							{
								WorkName: "operation",
							},
						},
					},
					Tests: []string{"1", "2"},
				},
				"name":  "张三",
				"age":   18,
				"citys": []string{"1", "2"},
				"ages":  []int64{1, 2},
			},
			Key:  "data.ages",
			Want: []int64{1, 2},
		},
		{
			Data: map[string]interface{}{
				"user": &Temp{
					User: &User{
						Name:   "zhangsan",
						Age:    18,
						IDCard: "110",
						Works: []*Work{
							{
								WorkName: "operation",
							},
						},
					},
					Tests: []string{"1", "2"},
				},
				"name":  "张三",
				"age":   18,
				"citys": []string{"1", "2"},
				"ages":  []int64{1, 2},
			},
			Key:  "data.user.User.Name",
			Want: "zhangsan",
		},
		{
			Data: map[string]interface{}{
				"user": &Temp{
					User: &User{
						Name:   "zhangsan",
						Age:    18,
						IDCard: "110",
						Works: []*Work{
							{
								WorkName: "operation",
							},
						},
					},
					Tests: []string{"1", "2"},
				},
				"name":  "张三",
				"age":   18,
				"citys": []string{"1", "2"},
				"ages":  []int64{1, 2},
			},
			Key: "data.user",
			Want: &Temp{
				User: &User{
					Name:   "zhangsan",
					Age:    18,
					IDCard: "110",
					Works: []*Work{
						{
							WorkName: "operation",
						},
					},
				},
				Tests: []string{"1", "2"},
			},
		},
		{
			Data: map[string]interface{}{
				"user": &Temp{
					User: &User{
						Name:   "zhangsan",
						Age:    18,
						IDCard: "110",
						Works: []*Work{
							{
								WorkName: "operation",
							},
						},
					},
					Tests: []string{"1", "2"},
				},
				"name":  "张三",
				"age":   18,
				"citys": []string{"1", "2"},
				"ages":  []int64{1, 2},
			},
			Key:  "data.user.Tests",
			Want: []string{"1", "2"},
		},
		{
			Data: map[string]interface{}{
				"user": &Temp{
					User: &User{
						Name:   "zhangsan",
						Age:    18,
						IDCard: "110",
						Works: []*Work{
							{
								WorkName: "operation",
							},
						},
					},
					Tests: []string{"1", "2"},
				},
				"name":  "张三",
				"age":   18,
				"citys": []string{"1", "2"},
				"ages":  []int64{1, 2},
			},
			Key:  "data.user.User.Age",
			Want: int8(18),
		},
		{
			Data: map[string]interface{}{
				"user": &Temp{
					User: &User{
						Name:   "zhangsan",
						Age:    18,
						IDCard: "110",
						Works: []*Work{
							{
								WorkName: "operation",
							},
						},
					},
					Tests: []string{"1", "2"},
				},
				"name":  "张三",
				"age":   18,
				"citys": []string{"1", "2"},
				"ages":  []int64{1, 2},
			},
			Key:  "data.user.User.Teacher",
			Want: nil,
			Err:  fmt.Errorf("%s not found in data", "data.user.User.Teacher"),
		},
		// array
		{
			Data: []interface{}{
				&Temp{
					User: &User{
						Name:   "zhangsan",
						Age:    18,
						IDCard: "110",
						Works: []*Work{
							{
								WorkName: "operation",
							},
						},
					},
					Tests: []string{"1", "2"},
				},
				1,
				"2",
				[]int{1, 2},
				[]string{"1", "2"},
			},
			Key:  "data.0.User.Name",
			Want: "zhangsan",
		},
		{
			Data: []interface{}{
				&Temp{
					User: &User{
						Name:   "zhangsan",
						Age:    18,
						IDCard: "110",
						Works: []*Work{
							{
								WorkName: "operation",
							},
						},
					},
					Tests: []string{"1", "2"},
				},
				1,
				"2",
				[]int{1, 2},
				[]string{"1", "2"},
			},
			Key: "data.0.User",
			Want: &User{
				Name:   "zhangsan",
				Age:    18,
				IDCard: "110",
				Works: []*Work{
					{
						WorkName: "operation",
					},
				},
			},
		},
		{
			Data: []interface{}{
				&Temp{
					User: &User{
						Name:   "zhangsan",
						Age:    18,
						IDCard: "110",
						Works: []*Work{
							{
								WorkName: "operation",
							},
						},
					},
					Tests: []string{"1", "2"},
				},
				1,
				"2",
				[]int{1, 2},
				[]string{"1", "2"},
			},
			Key:  "data.0.Tests",
			Want: []string{"1", "2"},
		},
		{
			Data: []interface{}{
				&Temp{
					User: &User{
						Name:   "zhangsan",
						Age:    18,
						IDCard: "110",
						Works: []*Work{
							{
								WorkName: "operation",
							},
						},
					},
					Tests: []string{"1", "2"},
				},
				1,
				"2",
				[]int{1, 2},
				[]string{"1", "2"},
			},
			Key:  "data.0.User.Age",
			Want: int8(18),
		},
		{
			Data: []interface{}{
				&Temp{
					User: &User{
						Name:   "zhangsan",
						Age:    18,
						IDCard: "110",
						Works: []*Work{
							{
								WorkName: "operation",
							},
						},
					},
					Tests: []string{"1", "2"},
				},
				1,
				"2",
				[]int{1, 2},
				[]string{"1", "2"},
			},
			Key:  "data.1",
			Want: 1,
		},
		{
			Data: []interface{}{
				&Temp{
					User: &User{
						Name:   "zhangsan",
						Age:    18,
						IDCard: "110",
						Works: []*Work{
							{
								WorkName: "operation",
							},
						},
					},
					Tests: []string{"1", "2"},
				},
				1,
				"2",
				[]int{1, 2},
				[]string{"1", "2"},
			},
			Key:  "data.2",
			Want: "2",
		},
		{
			Data: []interface{}{
				&Temp{
					User: &User{
						Name:   "zhangsan",
						Age:    18,
						IDCard: "110",
						Works: []*Work{
							{
								WorkName: "operation",
							},
						},
					},
					Tests: []string{"1", "2"},
				},
				1,
				"2",
				[]int{1, 2},
				[]string{"1", "2"},
			},
			Key:  "data.3",
			Want: []int{1, 2},
		},
		{
			Data: []interface{}{
				&Temp{
					User: &User{
						Name:   "zhangsan",
						Age:    18,
						IDCard: "110",
						Works: []*Work{
							{
								WorkName: "operation",
							},
						},
					},
					Tests: []string{"1", "2"},
				},
				1,
				"2",
				[]int{1, 2},
				[]string{"1", "2"},
			},
			Key:  "data.4",
			Want: []string{"1", "2"},
		},
		{
			Data: []interface{}{
				&Temp{
					User: &User{
						Name:   "zhangsan",
						Age:    18,
						IDCard: "110",
						Works: []*Work{
							{
								WorkName: "operation",
							},
						},
					},
					Tests: []string{"1", "2"},
				},
				1,
				"2",
				[]int{1, 2},
				[]string{"1", "2"},
			},
			Key:  "data.5",
			Want: nil,
			Err:  fmt.Errorf("%s not found in data", "data.5"),
		},
		// struct
		{
			Data: &Temp{
				User: &User{
					Name:   "zhangsan",
					Age:    18,
					IDCard: "110",
					Works: []*Work{
						{
							WorkName: "operation",
						},
					},
				},
				Tests: []string{"1", "2"},
			},
			Key:  "data.User.Name",
			Want: "zhangsan",
		},
		{
			Data: &Temp{
				User: &User{
					Name:   "zhangsan",
					Age:    18,
					IDCard: "110",
					Works: []*Work{
						{
							WorkName: "operation",
						},
					},
				},
				Tests: []string{"1", "2"},
			},
			Key:  "data.User.Age",
			Want: int8(18),
		},
		{
			Data: &Temp{
				User: &User{
					Name:   "zhangsan",
					Age:    18,
					IDCard: "110",
					Works: []*Work{
						{
							WorkName: "operation",
						},
					},
				},
				Tests: []string{"1", "2"},
			},
			Key:  "data.User.IDCard",
			Want: "110",
		},
		{
			Data: &Temp{
				User: &User{
					Name:   "zhangsan",
					Age:    18,
					IDCard: "110",
					Works: []*Work{
						{
							WorkName: "operation",
						},
					},
				},
				Tests: []string{"1", "2"},
			},
			Key: "data.User",
			Want: &User{
				Name:   "zhangsan",
				Age:    18,
				IDCard: "110",
				Works: []*Work{
					{
						WorkName: "operation",
					},
				},
			},
		},
		{
			Data: &Temp{
				User: &User{
					Name:   "zhangsan",
					Age:    18,
					IDCard: "110",
					Works: []*Work{
						{
							WorkName: "operation",
						},
					},
				},
				Tests: []string{"1", "2"},
			},
			Key:  "data.User.Works.0.WorkName",
			Want: "operation",
		},
		{
			Data: &Temp{
				User: &User{
					Name:   "zhangsan",
					Age:    18,
					IDCard: "110",
					Works: []*Work{
						{
							WorkName: "operation",
						},
					},
				},
				Tests: []string{"1", "2"},
			},
			Key:  "data.Tests",
			Want: []string{"1", "2"},
		},
		{
			Data: &Temp{
				User: &User{
					Name:   "zhangsan",
					Age:    18,
					IDCard: "110",
					Works: []*Work{
						{
							WorkName: "operation",
						},
					},
				},
				Tests: []string{"1", "2"},
			},
			Key:  "data.Age",
			Want: nil,
			Err:  fmt.Errorf("%s not found in data", "data.Age"),
		},
	}

	for index, tt := range cases {
		variable, ok := variables.Get(tt.Key)
		assert.True(t, ok)
		assert.NotNil(t, variable)
		assert.Equal(t, tt.Key, variable.Name(), index)

		ret, err := variable.Value(ctx, tt.Data, cc)
		if err != nil {
			assert.Equal(t, tt.Err, err)
		} else {
			assert.Equal(t, tt.Want, ret, index)
		}
	}
}
