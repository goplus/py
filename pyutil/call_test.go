package pyutil

import (
	"testing"
	"github.com/qiniu/x/log"
	"github.com/qiniu/x/errors"
	"github.com/qiniu/py"
)

type moduleCase struct {
	exp string
	name, ret, tbl string
}

var g_moduleCases = []moduleCase{
	{`
class Plugin:
	def init(self, cate):
		self.tbl = "dn_5m" + cate
`, "foo", "None", "dn_5m_stage"},
}

func TestCall(t *testing.T) {

	log.SetOutputLevel(0)

	for _, c := range g_moduleCases {
		code, err := py.Compile(c.exp, "", py.FileInput)
		if err != nil {
			t.Fatal("Compile failed:", err)
		}
		defer code.Decref()

		mod, err := py.ExecCodeModule(c.name, code.Obj())
		if err != nil {
			t.Fatal("ExecCodeModule failed:", err)
		}
		defer mod.Decref()

		plg, err := New(mod.Obj(), "Plugin")
		if err != nil {
			t.Fatal("NewPlugin failed:", errors.Detail(err))
		}

		ret, err := CallMethod(plg, "init", "_stage")
		if err != nil {
			t.Fatal("CallMethod failed:", err)
		}
		defer ret.Decref()

		if ret.String() != c.ret {
			t.Fatal("CallMethod ret:", ret.String(), c.ret)
		}

		tbl, _ := plg.GetAttrString("tbl")
		if tbl.String() != c.tbl {
			t.Fatal("mod.GetAttrString('tbl') ret:", tbl.String(), c.tbl)
		}
	}
}

