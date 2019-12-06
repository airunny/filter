package assignment

import (
	"reflect"
	"testing"
)

func TestEqual_Run(t *testing.T) {
	aa := make(map[string]interface{})

	reflect.ValueOf(aa).SetMapIndex(reflect.ValueOf("11"), reflect.ValueOf("11"))
	t.Log(aa)

	bb := []string{"11", "22", "33"}
	v := reflect.ValueOf(bb)
	e := v.Index(0)
	e.SetString("999")
	t.Log(bb)
}
