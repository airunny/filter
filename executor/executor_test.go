package executor

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/airunny/filter/assignment"
	_ "github.com/airunny/filter/assignment/set"
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

type mockAssignment struct {
	name string
	err  error
}

func (m mockAssignment) Name() string {
	return m.name
}

func (m mockAssignment) PrepareValue(ctx context.Context, value interface{}) (interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}
	return value, nil
}

func (m mockAssignment) Run(ctx context.Context, data interface{}, key string, val interface{}) error {
	return nil
}

func TestBuildExecutor(t *testing.T) {
	assignment.Register(&mockAssignment{
		name: "mock",
		err:  errors.New("mock error"),
	})

	cases := []struct {
		BuildErr error
		Data     interface{}
		Items    []interface{}
		Expected interface{}
	}{
		// err
		{
			BuildErr: errors.New("executor item must be array"),
			Items:    []interface{}{},
		},
		{
			BuildErr: errors.New("executor item must contains 3 elements"),
			Items:    []interface{}{"1", "2"},
		},
		{
			BuildErr: fmt.Errorf("executor item 1st item  %v is not string", 1),
			Items:    []interface{}{1, "2", "3"},
		},
		{
			BuildErr: fmt.Errorf("executor item 2nd item  %v is not string", 2),
			Items:    []interface{}{"1", 2, "3"},
		},
		{
			BuildErr: fmt.Errorf("executor assignment not exists [%s]", "append"),
			Items:    []interface{}{"name", "append", "3"},
		},
		{
			BuildErr: fmt.Errorf("executor assignment not exists [%s]", "append"),
			Items:    []interface{}{"name", "append", "3"},
		},
		{
			BuildErr: fmt.Errorf("executor assignment [%s] preparevalue err:%s", "mock", errors.New("mock error")),
			Items:    []interface{}{"name", "mock", "3"},
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
			Items: []interface{}{"name", "=", "李四"},
			Expected: map[string]interface{}{
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
				"name":  "李四",
				"age":   18,
				"citys": []string{"1", "2"},
				"ages":  []int64{1, 2},
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
			Items: []interface{}{"age", "=", 19},
			Expected: map[string]interface{}{
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
				"age":   19,
				"citys": []string{"1", "2"},
				"ages":  []int64{1, 2},
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
			Items: []interface{}{"citys", "=", []string{"3", "4"}},
			Expected: map[string]interface{}{
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
				"citys": []string{"3", "4"},
				"ages":  []int64{1, 2},
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
			Items: []interface{}{"ages", "=", []int64{3, 4}},
			Expected: map[string]interface{}{
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
				"ages":  []int64{3, 4},
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
			Items: []interface{}{"user.User.Name", "=", "zhangsan1"},
			Expected: map[string]interface{}{
				"user": &Temp{
					User: &User{
						Name:   "zhangsan1",
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
			Items: []interface{}{"user.User", "=", &User{
				Name:   "user",
				Age:    19,
				IDCard: "120",
				Works: []*Work{
					{
						WorkName: "operation1",
					},
				},
			}},
			Expected: map[string]interface{}{
				"user": &Temp{
					User: &User{
						Name:   "user",
						Age:    19,
						IDCard: "120",
						Works: []*Work{
							{
								WorkName: "operation1",
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
			Items: []interface{}{"user.Tests", "=", []string{"3", "4"}},
			Expected: map[string]interface{}{
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
					Tests: []string{"3", "4"},
				},
				"name":  "张三",
				"age":   18,
				"citys": []string{"1", "2"},
				"ages":  []int64{1, 2},
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
			Items: []interface{}{"user.User.Age", "=", 29},
			Expected: map[string]interface{}{
				"user": &Temp{
					User: &User{
						Name:   "zhangsan",
						Age:    29,
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
			Items: []interface{}{"0.User.Name", "=", "zhangsan2"},
			Expected: []interface{}{
				&Temp{
					User: &User{
						Name:   "zhangsan2",
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
			Items: []interface{}{"0.User", "=", &User{
				Name:   "user2",
				Age:    20,
				IDCard: "1200",
				Works: []*Work{
					{
						WorkName: "operation2",
					},
				},
			}},
			Expected: []interface{}{
				&Temp{
					User: &User{
						Name:   "user2",
						Age:    20,
						IDCard: "1200",
						Works: []*Work{
							{
								WorkName: "operation2",
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
			Items: []interface{}{"0.Tests", "=", []string{"30", "40"}},
			Expected: []interface{}{
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
					Tests: []string{"30", "40"},
				},
				1,
				"2",
				[]int{1, 2},
				[]string{"1", "2"},
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
			Items: []interface{}{"0.User.Age", "=", 40},
			Expected: []interface{}{
				&Temp{
					User: &User{
						Name:   "zhangsan",
						Age:    40,
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
			Items: []interface{}{"1", "=", 10},
			Expected: []interface{}{
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
				10,
				"2",
				[]int{1, 2},
				[]string{"1", "2"},
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
			Items: []interface{}{"2", "=", "20"},
			Expected: []interface{}{
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
				"20",
				[]int{1, 2},
				[]string{"1", "2"},
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
			Items: []interface{}{"3", "=", []int{3, 4}},
			Expected: []interface{}{
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
				[]int{3, 4},
				[]string{"1", "2"},
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
			Items: []interface{}{"4", "=", []string{"3", "4"}},
			Expected: []interface{}{
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
				[]string{"3", "4"},
			},
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
			Items: []interface{}{"User.Name", "=", "golang"},
			Expected: &Temp{
				User: &User{
					Name:   "golang",
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
			Items: []interface{}{"User.Age", "=", 17},
			Expected: &Temp{
				User: &User{
					Name:   "zhangsan",
					Age:    17,
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
			Items: []interface{}{"User.IDCard", "=", "id_card"},
			Expected: &Temp{
				User: &User{
					Name:   "zhangsan",
					Age:    18,
					IDCard: "id_card",
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
			Items: []interface{}{"User", "=", &User{
				Name:   "golang1",
				Age:    10,
				IDCard: "2009",
				Works: []*Work{
					{
						WorkName: "shanghai",
					},
				},
			}},
			Expected: &Temp{
				User: &User{
					Name:   "golang1",
					Age:    10,
					IDCard: "2009",
					Works: []*Work{
						{
							WorkName: "shanghai",
						},
					},
				},
				Tests: []string{"1", "2"},
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
			Items: []interface{}{"User.Works.0.WorkName", "=", "suzhou"},
			Expected: &Temp{
				User: &User{
					Name:   "zhangsan",
					Age:    18,
					IDCard: "110",
					Works: []*Work{
						{
							WorkName: "suzhou",
						},
					},
				},
				Tests: []string{"1", "2"},
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
			Items: []interface{}{"Tests", "=", []string{"3", "4"}},
			Expected: &Temp{
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
				Tests: []string{"3", "4"},
			},
		},
	}

	ctx := context.Background()
	for index, v := range cases {
		executor, err := BuildExecutor(ctx, v.Items)
		if err != nil {
			assert.True(t, reflect.DeepEqual(v.BuildErr, err), index)
		} else {
			err = executor.Execute(ctx, v.Data)
			assert.Nil(t, err)
			assert.True(t, reflect.DeepEqual(v.Expected, v.Data), index)
		}
	}
}
