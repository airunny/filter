package filter_type

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMyType(t *testing.T) {
	cases := []struct {
		Data     interface{}
		Expected FilterType
	}{
		// string
		{
			Data:     "",
			Expected: STRING,
		},
		{
			Data:     "111",
			Expected: STRING,
		},
		// number
		{
			Data:     int(0),
			Expected: NUMBER,
		},
		{
			Data:     int8(0),
			Expected: NUMBER,
		},
		{
			Data:     int16(0),
			Expected: NUMBER,
		},
		{
			Data:     int32(0),
			Expected: NUMBER,
		},
		{
			Data:     int64(0),
			Expected: NUMBER,
		},
		{
			Data:     uint(0),
			Expected: NUMBER,
		},
		{
			Data:     uint8(0),
			Expected: NUMBER,
		},
		{
			Data:     uint16(0),
			Expected: NUMBER,
		},
		{
			Data:     uint32(0),
			Expected: NUMBER,
		},
		{
			Data:     uint64(0),
			Expected: NUMBER,
		},
		{
			Data:     float32(0),
			Expected: NUMBER,
		},
		{
			Data:     float64(0),
			Expected: NUMBER,
		},
		// bool
		{
			Data:     true,
			Expected: BOOL,
		},
		{
			Data:     false,
			Expected: BOOL,
		},
		// hash
		{
			Data:     map[string]struct{}{},
			Expected: HASH,
		},
		{
			Data:     &map[string]struct{}{},
			Expected: HASH,
		},
		// struct
		{
			Data:     struct{}{},
			Expected: STRUCT,
		},
		{
			Data:     &struct{}{},
			Expected: STRUCT,
		},
		// array
		{
			Data:     []string{},
			Expected: ARRAY,
		},
		{
			Data:     &[]string{},
			Expected: ARRAY,
		},
		{
			Data:     [3]string{},
			Expected: ARRAY,
		},
		{
			Data:     &[3]string{},
			Expected: ARRAY,
		},
		// null
		{
			Data:     nil,
			Expected: NULL,
		},
	}

	for index, v := range cases {
		assert.Equal(t, v.Expected, GetFilterType(v.Data), index)
	}
}

type Student struct {
	Name string
	Age  int
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

func TestClone(t *testing.T) {
	stu := Student{
		Name: "zhangsan",
		Age:  18,
	}

	assert.Equal(t, true, reflect.DeepEqual(stu, Clone(stu)))
}
