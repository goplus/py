package py

import (
	"github.com/qiniu/x/log"
	"testing"
)

func init() {
	log.SetOutputLevel(0)
}

// ------------------------------------------------------------------------------------------

type Foo struct {
}

func (r *Foo) Py_foo(args *Tuple) (*Base, error) {
	return IncNone(), nil
}

func (r *Foo) Py_bar(args *Tuple) *Base {
	return IncNone()
}

// ------------------------------------------------------------------------------------------

func _TestRegister(t *testing.T) {

	dict := NewDict()
	defer dict.Decref()

	Register(dict, "", new(Foo))
}

// ------------------------------------------------------------------------------------------
