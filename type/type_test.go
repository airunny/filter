package _type

import (
	"fmt"
	"reflect"
	"testing"
)

type Student struct {
	Name string
	Age  int
}

func TestClone(t *testing.T) {
	stu := Student{
		Name: "zhangsan",
		Age:  18,
	}

	fmt.Println(reflect.TypeOf(stu).Kind())
}
