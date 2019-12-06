package _type

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMyType(t *testing.T) {
	cases := []struct {
		Data     interface{}
		Expected MyType
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
		assert.Equal(t, v.Expected, GetMyType(v.Data), index)
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
