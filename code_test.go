package py

import (
	"testing"
)

type compileCase struct {
	exp string
	ret string
	start StartToken
}

var g_compileCases = []compileCase{
	{"1+2", "3", EvalInput},
	{"1+2", "None", SingleInput}, // echo
	{"1+2", "None", FileInput},
}

func TestCompile(t *testing.T) {

	for _, c := range g_compileCases {
		code, err := Compile(c.exp, "", c.start)
		if err != nil {
			t.Fatal("Compile failed:", err)
		}
		defer code.Decref()

		globals := NewDict()
		defer globals.Decref()

		locals := NewDict()
		defer locals.Decref()

		ret, err := code.Eval(globals.Obj(), locals.Obj())
		if err != nil {
			t.Fatal("Eval failed:", err)
		}
		defer ret.Decref()

		if ret.String() != c.ret {
			t.Fatal("Eval ret:", ret.String())
		}
	}
}

type evalLocalGlobalsCase struct {
	exp string
	globals string
	locals string
	start StartToken
}

var g_evalLocalGlobalsCases = []evalLocalGlobalsCase{
	{"v=1+2", "{}", "{'v': 3}", FileInput},
	{"v=1+2", "{}", "{'v': 3}", SingleInput}, // echo
//	{"v=1+2", "{}", "{'v': 3}", EvalInput}, // compile error
}

func _TestEvalLocalGlobals(t *testing.T) {

	Initialize()
	defer Finalize()

	for _, c := range g_evalLocalGlobalsCases {
		code, err := Compile(c.exp, "", c.start)
		if err != nil {
			t.Fatal("Compile failed:", c.exp, c.start, err)
		}
		defer code.Decref()

		globals := NewDict()
		defer globals.Decref()

		locals := NewDict()
		defer locals.Decref()

		err = code.Run(globals.Obj(), locals.Obj())
		if err != nil {
			t.Fatal("Run failed:", err)
		}
		println(globals.String(), locals.String())

		if locals.String() != c.locals || globals.String() != c.globals {
			t.Fatal("Run:", globals.String(), locals.String())
		}
	}
}

