package py

import (
	"testing"
)

type moduleCase struct {
	exp                string
	name, ret, globals string
}

var g_moduleCases = []moduleCase{
	{`
tbl = 'dn_5m'
def init(cat):
	global tbl
	tbl = tbl + cat
	return True
	`, "foo", "True", "dn_5m_stage"},
}

func TestModule(t *testing.T) {

	for _, c := range g_moduleCases {
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

		arg1 := NewString("_stage")
		defer arg1.Decref()

		ret, err := mod.CallMethodObjArgs("init", arg1.Obj())
		if err != nil {
			t.Fatal("CallMethodObjArgs failed:", err)
		}
		defer ret.Decref()

		if ret.String() != c.ret {
			t.Fatal("CallMethodObjArgs ret:", ret.String(), c.ret)
		}

		globals, _ := mod.GetAttrString("tbl")
		defer globals.Decref()

		if globals.String() != c.globals {
			t.Fatal("mod.GetAttrString('tbl') ret:", globals.String(), c.globals)
		}

		dict := mod.Dict()                // don't need Decref
		tbl2 := dict.GetItemString("tbl") // don't need Decref
		if tbl2.String() != c.globals {
			t.Fatal("mod.Dict.GetItemString('tbl') ret:", tbl2.String(), c.globals)
		}
	}
}
