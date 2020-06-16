package executor

import (
	"context"
	"testing"

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

func TestBuildExecutor(t *testing.T) {
	structData := &Temp{
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
	}

	mapData := map[string]interface{}{"user": structData, "name": "张三", "age": 18, "citys": []string{"1", "2"}, "ages": []int64{1, 2}}
	arrayData := []interface{}{structData, 1, "2", []int{1, 2}, []string{"1", "2"}}

	cases := []struct {
		Data     interface{}
		Items    []interface{}
		Expected interface{}
		GotFunc  func() interface{}
	}{
		// map
		{
			Data:     mapData,
			Items:    []interface{}{"name", "=", "李四"},
			Expected: "李四",
			GotFunc: func() interface{} {
				return mapData["name"]
			},
		},
		{
			Data:     mapData,
			Items:    []interface{}{"age", "=", 19},
			Expected: 19,
			GotFunc: func() interface{} {
				return mapData["age"]
			},
		},
		{
			Data:     mapData,
			Items:    []interface{}{"citys", "=", []string{"3", "4"}},
			Expected: []string{"3", "4"},
			GotFunc: func() interface{} {
				return mapData["citys"]
			},
		},
		{
			Data:     mapData,
			Items:    []interface{}{"ages", "=", []int64{3, 4}},
			Expected: []int64{3, 4},
			GotFunc: func() interface{} {
				return mapData["ages"]
			},
		},
		{
			Data:     mapData,
			Items:    []interface{}{"user.User.Name", "=", "zhangsan1"},
			Expected: "zhangsan1",
			GotFunc: func() interface{} {
				return structData.User.Name
			},
		},
		{
			Data: mapData,
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
			Expected: &User{
				Name:   "user",
				Age:    19,
				IDCard: "120",
				Works: []*Work{
					{
						WorkName: "operation1",
					},
				},
			},
			GotFunc: func() interface{} {
				return structData.User
			},
		},
		{
			Data:     mapData,
			Items:    []interface{}{"user.Tests", "=", []string{"3", "4"}},
			Expected: []string{"3", "4"},
			GotFunc: func() interface{} {
				return structData.Tests
			},
		},
		{
			Data:     mapData,
			Items:    []interface{}{"user.User.Age", "=", 29},
			Expected: int8(29),
			GotFunc: func() interface{} {
				return structData.User.Age
			},
		},

		// array
		{
			Data:     arrayData,
			Items:    []interface{}{"0.User.Name", "=", "zhangsan2"},
			Expected: "zhangsan2",
			GotFunc: func() interface{} {
				return structData.User.Name
			},
		},
		{
			Data: arrayData,
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
			Expected: &User{
				Name:   "user2",
				Age:    20,
				IDCard: "1200",
				Works: []*Work{
					{
						WorkName: "operation2",
					},
				},
			},
			GotFunc: func() interface{} {
				return structData.User
			},
		},
		{
			Data:     arrayData,
			Items:    []interface{}{"0.Tests", "=", []string{"30", "40"}},
			Expected: []string{"30", "40"},
			GotFunc: func() interface{} {
				return structData.Tests
			},
		},
		{
			Data:     arrayData,
			Items:    []interface{}{"0.User.Age", "=", 40},
			Expected: int8(40),
			GotFunc: func() interface{} {
				return structData.User.Age
			},
		},
		{
			Data:     arrayData,
			Items:    []interface{}{"1", "=", 10},
			Expected: 10,
			GotFunc: func() interface{} {
				return arrayData[1]
			},
		},
		{
			Data:     arrayData,
			Items:    []interface{}{"2", "=", "20"},
			Expected: "20",
			GotFunc: func() interface{} {
				return arrayData[2]
			},
		},
		{
			Data:     arrayData,
			Items:    []interface{}{"3", "=", []int{3, 4}},
			Expected: []int{3, 4},
			GotFunc: func() interface{} {
				return arrayData[3]
			},
		},
		{
			Data:     arrayData,
			Items:    []interface{}{"4", "=", []string{"3", "4"}},
			Expected: []string{"3", "4"},
			GotFunc: func() interface{} {
				return arrayData[4]
			},
		},
		// struct
		{
			Data:     structData,
			Items:    []interface{}{"User.Name", "=", "golang"},
			Expected: "golang",
			GotFunc: func() interface{} {
				return structData.User.Name
			},
		},
		{
			Data:     structData,
			Items:    []interface{}{"User.Age", "=", 17},
			Expected: int8(17),
			GotFunc: func() interface{} {
				return structData.User.Age
			},
		},
		{
			Data:     structData,
			Items:    []interface{}{"User.IDCard", "=", "id_card"},
			Expected: "id_card",
			GotFunc: func() interface{} {
				return structData.User.IDCard
			},
		},
		{
			Data: structData,
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
			Expected: &User{
				Name:   "golang1",
				Age:    10,
				IDCard: "2009",
				Works: []*Work{
					{
						WorkName: "shanghai",
					},
				},
			},
			GotFunc: func() interface{} {
				return structData.User
			},
		},
		{
			Data:     structData,
			Items:    []interface{}{"User.Works.0.WorkName", "=", "suzhou"},
			Expected: "suzhou",
			GotFunc: func() interface{} {
				return structData.User.Works[0].WorkName
			},
		},
		{
			Data:     structData,
			Items:    []interface{}{"Tests", "=", []string{"3", "4"}},
			Expected: []string{"3", "4"},
			GotFunc: func() interface{} {
				return structData.Tests
			},
		},
	}

	ctx := context.Background()
	for index, v := range cases {
		executor, err := BuildExecutor(ctx, v.Items)
		assert.Equal(t, nil, err, index)
		executor.Execute(ctx, v.Data)

		assert.Equal(t, v.Expected, v.GotFunc(), index)
	}
}

func TestBuildExecutor2(t *testing.T) {
	structData := &Temp{
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
	}

	mapData := map[string]interface{}{"user": structData, "name": "张三", "age": 18, "citys": []string{"1", "2"}, "ages": []int64{1, 2}}
	arrayData := []interface{}{structData, 1, "2", []int{1, 2}, []string{"1", "2"}}

	cases := []struct {
		Data     interface{}
		Items    []interface{}
		Expected interface{}
		GotFunc  func() interface{}
	}{
		// map
		{
			Data: mapData,
			Items: []interface{}{
				[]interface{}{"name", "=", "李四"},
				[]interface{}{"age", "=", 19},
				[]interface{}{"citys", "=", []string{"3", "4"}},
				[]interface{}{"ages", "=", []int64{3, 4}},
				[]interface{}{"user.User", "=", &User{
					Name:   "user",
					Age:    19,
					IDCard: "120",
					Works: []*Work{
						{
							WorkName: "operation1",
						},
					},
				}},
				[]interface{}{"user.Tests", "=", []string{"3", "4"}},
				[]interface{}{"user.User.Age", "=", 29},
				[]interface{}{"user.User.Name", "=", "zhangsan1"},
			},
			Expected: map[string]interface{}{"user": &Temp{
				User: &User{
					Name:   "zhangsan1",
					Age:    29,
					IDCard: "120",
					Works: []*Work{
						{
							WorkName: "operation1",
						},
					},
				},
				Tests: []string{"3", "4"},
			}, "name": "李四", "age": 19, "citys": []string{"3", "4"}, "ages": []int64{3, 4}},
			GotFunc: func() interface{} {
				return mapData
			},
		},

		// array
		{
			Data: arrayData,
			Items: []interface{}{
				[]interface{}{"0.User", "=", &User{
					Name:   "user2",
					Age:    20,
					IDCard: "1200",
					Works: []*Work{
						{
							WorkName: "operation2",
						},
					},
				}},
				[]interface{}{"0.User.Name", "=", "zhangsan2"},
				[]interface{}{"0.User.Age", "=", 40},
				[]interface{}{"0.User.IDCard", "=", "id_card"},
				[]interface{}{"0.Tests", "=", []string{"30", "40"}},
				[]interface{}{"1", "=", 10},
				[]interface{}{"2", "=", "20"},
				[]interface{}{"3", "=", []int{3, 4}},
				[]interface{}{"4", "=", []string{"3", "4"}},
			},
			Expected: []interface{}{&Temp{
				User: &User{
					Name:   "zhangsan2",
					Age:    40,
					IDCard: "id_card",
					Works: []*Work{
						{
							WorkName: "operation2",
						},
					},
				},
				Tests: []string{"30", "40"},
			}, 10, "20", []int{3, 4}, []string{"3", "4"}},
			GotFunc: func() interface{} {
				return arrayData
			},
		},

		// struct
		{
			Data: structData,
			Items: []interface{}{
				[]interface{}{"User", "=", &User{
					Name:   "golang1",
					Age:    10,
					IDCard: "2009",
					Works: []*Work{
						{
							WorkName: "shanghai",
						},
					},
				}},
				[]interface{}{"User.Name", "=", "golang"},
				[]interface{}{"User.Age", "=", 17},
				[]interface{}{"User.IDCard", "=", "id_card"},
				[]interface{}{"User.Works.0.WorkName", "=", "suzhou"},
				[]interface{}{"Tests", "=", []string{"3", "4"}},
			},
			Expected: &Temp{
				User: &User{
					Name:   "golang",
					Age:    17,
					IDCard: "id_card",
					Works: []*Work{
						{
							WorkName: "suzhou",
						},
					},
				},
				Tests: []string{"3", "4"},
			},
			GotFunc: func() interface{} {
				return structData
			},
		},
	}

	ctx := context.Background()
	for index, v := range cases {
		executor, err := BuildExecutor(ctx, v.Items)
		assert.Equal(t, nil, err, index)
		executor.Execute(ctx, v.Data)

		assert.Equal(t, v.Expected, v.GotFunc(), index)
	}
}

func TestBuildExecutor3(t *testing.T) {
	structData := &Temp{
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
	}

	mapData := map[string]interface{}{"user": structData, "name": "张三", "age": 18, "citys": []string{"1", "2"}, "ages": []int64{1, 2}}
	arrayData := []interface{}{structData, 1, "2", []int{1, 2}, []string{"1", "2"}}

	cases := []struct {
		Data     interface{}
		Items    []interface{}
		Expected interface{}
		GotFunc  func() interface{}
	}{
		// map
		{
			Data: mapData,
			Items: []interface{}{
				"__set", "=", []interface{}{
					[]interface{}{"name", "=", "李四"},
					[]interface{}{"age", "=", 19},
					[]interface{}{"citys", "=", []string{"3", "4"}},
					[]interface{}{"ages", "=", []int64{3, 4}},
					[]interface{}{"user.User", "=", &User{
						Name:   "user",
						Age:    19,
						IDCard: "120",
						Works: []*Work{
							{
								WorkName: "operation1",
							},
						},
					}},
					[]interface{}{"user.Tests", "=", []string{"3", "4"}},
					[]interface{}{"user.User.Age", "=", 29},
					[]interface{}{"user.User.Name", "=", "zhangsan1"},
				},
			},
			Expected: map[string]interface{}{"user": &Temp{
				User: &User{
					Name:   "zhangsan1",
					Age:    29,
					IDCard: "120",
					Works: []*Work{
						{
							WorkName: "operation1",
						},
					},
				},
				Tests: []string{"3", "4"},
			}, "name": "李四", "age": 19, "citys": []string{"3", "4"}, "ages": []int64{3, 4}},
			GotFunc: func() interface{} {
				return mapData
			},
		},

		// array
		{
			Data: arrayData,
			Items: []interface{}{
				"__set", "=", []interface{}{
					[]interface{}{"0.User", "=", &User{
						Name:   "user2",
						Age:    20,
						IDCard: "1200",
						Works: []*Work{
							{
								WorkName: "operation2",
							},
						},
					}},
					[]interface{}{"0.User.Name", "=", "zhangsan2"},
					[]interface{}{"0.User.Age", "=", 40},
					[]interface{}{"0.User.IDCard", "=", "id_card"},
					[]interface{}{"0.Tests", "=", []string{"30", "40"}},
					[]interface{}{"1", "=", 10},
					[]interface{}{"2", "=", "20"},
					[]interface{}{"3", "=", []int{3, 4}},
					[]interface{}{"4", "=", []string{"3", "4"}},
				},
			},
			Expected: []interface{}{&Temp{
				User: &User{
					Name:   "zhangsan2",
					Age:    40,
					IDCard: "id_card",
					Works: []*Work{
						{
							WorkName: "operation2",
						},
					},
				},
				Tests: []string{"30", "40"},
			}, 10, "20", []int{3, 4}, []string{"3", "4"}},
			GotFunc: func() interface{} {
				return arrayData
			},
		},

		// struct
		{
			Data: structData,
			Items: []interface{}{
				"__set", "=", []interface{}{
					[]interface{}{"User", "=", &User{
						Name:   "golang1",
						Age:    10,
						IDCard: "2009",
						Works: []*Work{
							{
								WorkName: "shanghai",
							},
						},
					}},
					[]interface{}{"User.Name", "=", "golang"},
					[]interface{}{"User.Age", "=", 17},
					[]interface{}{"User.IDCard", "=", "id_card"},
					[]interface{}{"User.Works.0.WorkName", "=", "suzhou"},
					[]interface{}{"Tests", "=", []string{"3", "4"}},
				},
			},
			Expected: &Temp{
				User: &User{
					Name:   "golang",
					Age:    17,
					IDCard: "id_card",
					Works: []*Work{
						{
							WorkName: "suzhou",
						},
					},
				},
				Tests: []string{"3", "4"},
			},
			GotFunc: func() interface{} {
				return structData
			},
		},
	}

	ctx := context.Background()
	for index, v := range cases {
		executor, err := BuildExecutor(ctx, v.Items)
		assert.Equal(t, nil, err, index)
		executor.Execute(ctx, v.Data)

		assert.Equal(t, v.Expected, v.GotFunc(), index)
	}
}
