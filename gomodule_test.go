package py

import (
	"fmt"
	"testing"
)

// ------------------------------------------------------------------------------------------

type FooModule struct {
}

func (r *FooModule) Py_bar(args *Tuple) (ret *Base, err error) {
	var i int
	var s []string
	err = ParseV(args, &i, &s)
	if err != nil {
		return
	}
	fmt.Println("call foo.bar:", i, s)
	return IncNone(), nil
}

// ------------------------------------------------------------------------------------------

type gomoduleCase struct {
	exp string
	name string
}

var g_gomoduleCases = []gomoduleCase{
	{
`import foo
foo.bar(1, 'Hello')
`, "test"},
}

func TestGoModule(t *testing.T) {

	gomod, err := NewGoModule("foo", "", new(FooModule))
	if err != nil {
		t.Fatal("NewGoModule failed:", err)
	}
	defer gomod.Decref()

	for _, c := range g_gomoduleCases {

		code, err := Compile(c.exp, "", FileInput)
		if err != nil {
			t.Fatal("Compile failed:", err)
		}
		defer code.Decref()

		mod, err := ExecCodeModule(c.name, code.Obj())
		if err != nil {
			t.Fatal("ExecCodeModule failed:", err)
		}
		defer mod.Decref()
	}
}

// ------------------------------------------------------------------------------------------

