package py

import (
	"testing"
)

func TestBase(t *testing.T) {

	{
		v := NewString("Hello!")
		defer v.Decref()

		if v.String() != "Hello!" {
			t.Fatal("NewString failed")
		}
	}
	{
		v := NewInt(1)
		defer v.Decref()

		if v.String() != "1" {
			t.Fatal("NewInt failed")
		}
	}
	{
		v1 := NewInt(1)
		defer v1.Decref()

		v2 := NewString("Hello!")
		defer v2.Decref()

		v := PackTuple(v1.Obj(), v2.Obj())
		defer v.Decref()

		if v.String() != "(1, 'Hello!')" {
			t.Fatal("NewTuple failed:", v.String())
		}
	}
	{
		v1 := NewInt(1)
		defer v1.Decref()

		v2 := NewString("Hello!")
		defer v2.Decref()

		v := NewDict()
		defer v.Decref()

		v.SetItem(v1.Obj(), v2.Obj())

		if v.String() != "{1: 'Hello!'}" {
			t.Fatal("NewDict failed:", v.String())
		}
	}
}
