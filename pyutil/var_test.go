package pyutil

import (
	"testing"
	"github.com/qiniu/py"
)

type Foo struct {
	A	int		`json:"a"`
	B	string	`json:"b"`
}

func Test(t *testing.T) {

	{
		val, ok := NewVar(1)
		if !ok {
			t.Fatal("NewVar failed")
		}
		if v, ok := py.AsInt(val); !ok || v.Int() != 1 {
			t.Fatal("NewVar failed:", val)
		}
	}
	{
		val, ok := NewVar(int64(1))
		if !ok {
			t.Fatal("NewVar failed")
		}
		if v, ok := py.AsLong(val); !ok || v.Long() != 1 {
			t.Fatal("NewVar failed:", val)
		}
	}
	{
		val, ok := NewVar("Hello")
		if !ok {
			t.Fatal("NewVar failed")
		}
		if v, ok := py.AsString(val); !ok || v.String() != "Hello" {
			t.Fatal("NewVar failed:", val)
		}
	}
	{
		foo := &Foo{
			A: 1, B: "Hello",
		}
		val, ok := NewVar(foo)
		if !ok {
			t.Fatal("NewVar failed")
		}
		if v, ok := py.AsDict(val); !ok || !checkFoo(v, t) {
			t.Fatal("NewVar failed:", val)
		}
	}
	{
		foo := map[string]interface{}{
			"a": 1, "b": "Hello",
		}
		val, ok := NewVar(foo)
		if !ok {
			t.Fatal("NewVar failed")
		}
		if v, ok := py.AsDict(val); !ok || !checkFoo(v, t) {
			t.Fatal("NewVar failed:", val)
		}
	}
}

func checkFoo(val *py.Dict, t *testing.T) bool {

	a := val.GetItemString("a")
	if a == nil {
		t.Fatal("GetItemString a failed")
		return false
	}
	if v, ok := py.AsInt(a); !ok || v.Int() != 1 {
		t.Fatal("GetItemString a failed")
	}

	b := val.GetItemString("b")
	if b == nil {
		t.Fatal("GetItemString b failed")
		return false
	}
	if v, ok := py.AsString(b); !ok || v.String() != "Hello" {
		t.Fatal("GetItemString b failed")
	}
	return true
}

