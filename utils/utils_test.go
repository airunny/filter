package utils

import (
	"math"
	"net"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloatEquals(t *testing.T) {
	cases := []struct {
		A        float64
		B        float64
		Expected bool
	}{
		{
			A:        0,
			B:        0 + EPSILON,
			Expected: false,
		},
		{
			A:        0,
			B:        0 + 0.000000001,
			Expected: true,
		},
		{
			A:        -1,
			B:        -1 + 0.000000001,
			Expected: true,
		},
		{
			A:        math.MaxFloat64,
			B:        math.MaxFloat64 - 0.0000001,
			Expected: true,
		},
		{
			A:        1,
			B:        1,
			Expected: true,
		},
	}

	for i, v := range cases {
		assert.Equal(t, v.Expected, FloatEquals(v.A, v.B), i)
	}
}

// -------------
type Work struct {
	Name string `json:"name"`
}

type User struct {
	Name   string `json:"name"`
	Age    int32  `json:"age"`
	IDCard string `json:"id_card"`
	Works  []Work `json:"works"`
}

type Temp struct {
	User User `json:"user"`
}

func TestGetObjectValueByKey(t *testing.T) {
	mock := Temp{
		User: User{
			Name:   "zhangsan",
			Age:    18,
			IDCard: "110",
			Works: []Work{
				{
					Name: "111",
				},
				{
					Name: "222",
				},
			},
		}}

	cases := []struct {
		Data     interface{}
		Key      string
		Expected interface{}
		OK       bool
	}{
		// map
		{
			Data:     map[string]interface{}{"user": "zhangsan", "age": 18},
			Key:      "user",
			Expected: "zhangsan",
			OK:       true,
		},
		{
			Data:     map[string]interface{}{"user": "zhangsan", "age": 18},
			Key:      "age",
			Expected: 18,
			OK:       true,
		},
		{
			Data:     map[string]interface{}{"user": mock},
			Key:      "user",
			Expected: mock,
			OK:       true,
		},
		{
			Data:     map[string]interface{}{"user": mock},
			Key:      "user.User",
			Expected: mock.User,
			OK:       true,
		},
		{
			Data:     map[string]interface{}{"user": mock},
			Key:      "user.User.Name",
			Expected: mock.User.Name,
			OK:       true,
		},
		{
			Data:     map[string]interface{}{"user": mock},
			Key:      "user.User.Age",
			Expected: mock.User.Age,
			OK:       true,
		},
		{
			Data:     map[string]interface{}{"user": mock},
			Key:      "user.User.Works.0",
			Expected: mock.User.Works[0],
			OK:       true,
		},
		{
			Data:     map[string]interface{}{"user": mock},
			Key:      "user.User.Works.0.Name",
			Expected: mock.User.Works[0].Name,
			OK:       true,
		},
		// map ptr
		{
			Data:     &map[string]interface{}{"user": "zhangsan", "age": 18},
			Key:      "user",
			Expected: "zhangsan",
			OK:       true,
		},
		{
			Data:     &map[string]interface{}{"user": "zhangsan", "age": 18},
			Key:      "age",
			Expected: 18,
			OK:       true,
		},
		{
			Data:     &map[string]interface{}{"user": mock},
			Key:      "user",
			Expected: mock,
			OK:       true,
		},
		{
			Data:     &map[string]interface{}{"user": mock},
			Key:      "user.User",
			Expected: mock.User,
			OK:       true,
		},
		{
			Data:     &map[string]interface{}{"user": mock},
			Key:      "user.User.Name",
			Expected: mock.User.Name,
			OK:       true,
		},
		{
			Data:     &map[string]interface{}{"user": mock},
			Key:      "user.User.Age",
			Expected: mock.User.Age,
			OK:       true,
		},
		{
			Data:     &map[string]interface{}{"user": mock},
			Key:      "user.User.Works.0",
			Expected: mock.User.Works[0],
			OK:       true,
		},
		{
			Data:     &map[string]interface{}{"user": mock},
			Key:      "user.User.Works.0.Name",
			Expected: mock.User.Works[0].Name,
			OK:       true,
		},
		// array
		{
			Data:     []interface{}{mock},
			Key:      "0",
			Expected: mock,
			OK:       true,
		},
		{
			Data:     []interface{}{mock},
			Key:      "0.User",
			Expected: mock.User,
			OK:       true,
		},
		{
			Data:     []interface{}{mock},
			Key:      "0.User.Name",
			Expected: mock.User.Name,
			OK:       true,
		},
		{
			Data:     []interface{}{mock},
			Key:      "0.User.Age",
			Expected: mock.User.Age,
			OK:       true,
		},
		{
			Data:     []interface{}{mock},
			Key:      "0.User.Works.0",
			Expected: mock.User.Works[0],
			OK:       true,
		},
		{
			Data:     []interface{}{mock},
			Key:      "0.User.Works.0.Name",
			Expected: mock.User.Works[0].Name,
			OK:       true,
		},
		// array ptr
		{
			Data:     &[]interface{}{mock},
			Key:      "0",
			Expected: mock,
			OK:       true,
		},
		{
			Data:     &[]interface{}{mock},
			Key:      "0.User",
			Expected: mock.User,
			OK:       true,
		},
		{
			Data:     &[]interface{}{mock},
			Key:      "0.User.Name",
			Expected: mock.User.Name,
			OK:       true,
		},
		{
			Data:     &[]interface{}{mock},
			Key:      "0.User.Age",
			Expected: mock.User.Age,
			OK:       true,
		},
		{
			Data:     &[]interface{}{mock},
			Key:      "0.User.Works.0",
			Expected: mock.User.Works[0],
			OK:       true,
		},
		{
			Data:     &[]interface{}{mock},
			Key:      "0.User.Works.0.Name",
			Expected: mock.User.Works[0].Name,
			OK:       true,
		},
		// struct
		{
			Data:     mock,
			Key:      ".",
			Expected: mock,
			OK:       true,
		},
		{
			Data:     mock,
			Key:      "",
			Expected: mock,
			OK:       true,
		},
		{
			Data:     mock,
			Key:      "User",
			Expected: mock.User,
			OK:       true,
		},
		{
			Data:     mock,
			Key:      "User.Name",
			Expected: mock.User.Name,
			OK:       true,
		},
		{
			Data:     mock,
			Key:      "User.Age",
			Expected: mock.User.Age,
			OK:       true,
		},
		{
			Data:     mock,
			Key:      "User.Works.0",
			Expected: mock.User.Works[0],
			OK:       true,
		},
		{
			Data:     mock,
			Key:      "User.Works.0.Name",
			Expected: mock.User.Works[0].Name,
			OK:       true,
		},
		// struct ptr
		{
			Data:     &mock,
			Key:      ".",
			Expected: &mock,
			OK:       true,
		},
		{
			Data:     &mock,
			Key:      "",
			Expected: &mock,
			OK:       true,
		},
		{
			Data:     &mock,
			Key:      "User",
			Expected: mock.User,
			OK:       true,
		},
		{
			Data:     &mock,
			Key:      "User.Name",
			Expected: mock.User.Name,
			OK:       true,
		},
		{
			Data:     &mock,
			Key:      "User.Age",
			Expected: mock.User.Age,
			OK:       true,
		},
		{
			Data:     &mock,
			Key:      "User.Works.0",
			Expected: mock.User.Works[0],
			OK:       true,
		},
		{
			Data:     &mock,
			Key:      "User.Works.0.Name",
			Expected: mock.User.Works[0].Name,
			OK:       true,
		},
	}

	for index, v := range cases {
		ret, ok := GetObjectValueByKey(v.Data, v.Key)
		assert.Equal(t, v.OK, ok, index)
		assert.Equal(t, true, reflect.DeepEqual(v.Expected, ret))
	}
}

func TestObjectCompare(t *testing.T) {
	cases := []struct {
		compare  interface{}
		compared interface{}
		expected int
	}{
		// int
		{
			compare:  1,
			compared: 1,
			expected: 0,
		},
		{
			compare:  2,
			compared: 1,
			expected: 1,
		},
		{
			compare:  1,
			compared: 2,
			expected: -1,
		},
		// float
		{
			compare:  1.0,
			compared: 1.0,
			expected: 0,
		},
		{
			compare:  2.0,
			compared: 1.0,
			expected: 1,
		},
		{
			compare:  1.0,
			compared: 2.0,
			expected: -1,
		},
		// float and int
		{
			compare:  1,
			compared: 1.0,
			expected: 0,
		},
		{
			compare:  2,
			compared: 1.0,
			expected: 1,
		},
		{
			compare:  1,
			compared: 2.0,
			expected: -1,
		},
		// string
		{
			compare:  "1",
			compared: "1",
			expected: 0,
		},
		{
			compare:  "2",
			compared: "1",
			expected: 1,
		},
		{
			compare:  "1",
			compared: "2",
			expected: -1,
		},
		// bool
		{
			compare:  true,
			compared: true,
			expected: 0,
		},
		{
			compare:  1,
			compared: false,
			expected: 1,
		},
		{
			compare:  false,
			compared: 1.0,
			expected: -1,
		},
		// arr
		{
			compare:  []string{"1", "2"},
			compared: []string{"1", "2"},
			expected: 0,
		},
		{
			compare:  []string{"1", "2", "3"},
			compared: []string{"1", "2"},
			expected: 1,
		},
		{
			compare:  []int{1, 2},
			compared: []int{1, 2},
			expected: 0,
		},
		{
			compare:  []float64{1, 2},
			compared: []float64{1, 2},
			expected: 0,
		},
	}

	for i, v := range cases {
		assert.Equal(t, v.expected, ObjectCompare(v.compare, v.compared), i)
	}
}

type Student struct {
	Name string
	Age  int
}

func TestClone(t *testing.T) {
	stu := Student{
		Name: "zhangsan",
		Age:  18,
	}

	assert.Equal(t, true, reflect.DeepEqual(stu, Clone(stu)))
}

func TestVersionCompare(t *testing.T) {
	cases := []struct {
		Compare  string
		Compared string
		Result   int
	}{
		{
			Compare:  "",
			Compared: "",
			Result:   0,
		},
		{
			Compare:  "1",
			Compared: "",
			Result:   1,
		},
		{
			Compare:  "",
			Compared: "1",
			Result:   -1,
		},
		{
			Compare:  "1.1.1",
			Compared: "1.1.1",
			Result:   0,
		},
		{
			Compare:  "1.2.1",
			Compared: "1.1.1",
			Result:   1,
		},
		{
			Compare:  "1.1.1",
			Compared: "1.2.1",
			Result:   -1,
		},
		{
			Compare:  "1.2.1",
			Compared: "1.1.2",
			Result:   1,
		},
		{
			Compare:  "1.1.1.1",
			Compared: "1.2.1",
			Result:   -1,
		},
		{
			Compare:  "1.1.1.0",
			Compared: "1.1.1",
			Result:   0,
		},
		{
			Compare:  "1.1.1",
			Compared: "1.1.1.0",
			Result:   0,
		},
	}

	for index, v := range cases {
		assert.Equal(t, v.Result, VersionCompare(v.Compare, v.Compared), index)
	}
}

func TestToInt(t *testing.T) {
	ip := net.ParseIP("192.168.2.1")
	expected := 192<<24 + 168<<16 + 2<<8 + 1
	assert.Equal(t, expected, ToInt(ip))
}

func TestIntToIP(t *testing.T) {
	cases := []struct {
		IP string
	}{
		{
			IP: "0.0.0.0",
		},
		{
			IP: "192.168.2.1",
		},
		{
			IP: "255.255.255.255",
		},
	}

	for _, v := range cases {
		assert.Equal(t, v.IP, IntToIP(ToInt(net.ParseIP(v.IP))).String())
	}

}

func TestBytesNOT(t *testing.T) {
	cases := []struct {
		B        []byte
		Expected []byte
	}{
		{
			B:        []byte{0, 1, 2, 3, 4, 100, 200, 250, 255},
			Expected: []byte{255 - 0, 255 - 1, 255 - 2, 255 - 3, 255 - 4, 255 - 100, 255 - 200, 255 - 250, 255 - 255},
		},
		{
			B:        []byte{},
			Expected: []byte{},
		},
		{
			B:        []byte{255, 255, 255, 255, 255},
			Expected: []byte{0, 0, 0, 0, 0},
		},
		{
			B:        []byte{0, 0, 0, 0, 0},
			Expected: []byte{255, 255, 255, 255, 255},
		},
	}

	for _, v := range cases {
		got := BytesNOT(v.B)
		if !reflect.DeepEqual(v.Expected, got) {
			t.Errorf("expected: %v,but Got:%v\n", v.Expected, got)
		}
	}
}

func TestBytesOR(t *testing.T) {
	cases := []struct {
		A []byte
		B []byte
	}{
		{
			A: []byte{1, 11, 20, 33, 50, 100, 200, 255, 254},
			B: []byte{0, 255, 0, 100, 0, 255, 199, 233, 255},
		},
		{
			A: []byte{},
			B: []byte{},
		},
		{
			A: []byte{255, 255, 255, 255, 255},
			B: []byte{255, 255, 255, 255, 255},
		},
		{
			A: []byte{0, 0, 0, 0, 0},
			B: []byte{0, 0, 0, 0, 0},
		},
		{
			A: []byte{255, 255, 255, 255, 255},
			B: []byte{0, 0, 0, 0, 0},
		},
	}

	for index, v := range cases {
		got := BytesOR(v.A, v.B)
		for j := 0; j < len(got); j++ {
			assert.Equal(t, v.A[j]|v.B[j], got[j], index)
		}
	}
}
