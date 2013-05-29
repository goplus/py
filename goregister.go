package py

import "reflect"
import "strings"
import "github.com/qiniu/log"

// ------------------------------------------------------------------------------------------

// Note: Methods take the receiver as the first argument, which the want
// signature doesn't include.
func sigMatches(got, want reflect.Type) bool {

	nin := want.NumIn()
	if got.NumIn()-1 != nin {
		return false
	}

	nout := want.NumOut()
	if got.NumOut() != nout {
		return false
	}

	for i := 0; i < nin; i++ {
		if got.In(i+1) != want.In(i) {
			return false
		}
	}

	for i := 0; i < nout; i++ {
		if got.Out(i) != want.Out(i) {
			return false
		}
	}
	return true
}

// ------------------------------------------------------------------------------------------

var typUnaryFunc = reflect.TypeOf(func() (*Base, error)(nil))
var typBinaryCallFunc = reflect.TypeOf(func(*Tuple) (*Base, error)(nil))
var typTernaryCallFunc = reflect.TypeOf(func(*Tuple, *Dict) (*Base, error)(nil))

type RegisterCtx []*Closure // 只是让对象不被gc

func Register(dict *Dict, nsprefix string, self interface{}) (ctx RegisterCtx) {

	typ := reflect.TypeOf(self)
	selfv := reflect.ValueOf(self)

	nmethod := typ.NumMethod()

	for i := 0; i < nmethod; i++ {
		method := typ.Method(i)
		mtype := method.Type
		mname := method.Name
		if mtype.PkgPath() != "" || !strings.HasPrefix(mname, "Py_") {
			continue
		}
		nin := mtype.NumIn()
		name := mname[3:]
		fullname := nsprefix + name
		if nin == 3 && sigMatches(mtype, typTernaryCallFunc) || nin == 2 && sigMatches(mtype, typBinaryCallFunc) {
			closure := &Closure{Self: selfv, Method: method.Func}
			f := closure.NewFunction(fullname, nin, "")
			dict.SetItemString(name, f)
			f.Decref()
			ctx = append(ctx, closure)
			log.Debug("Register", fullname)
		} else {
			log.Warnf("Invalid signature of method %s, register failed", fullname)
			continue
		}
	}
	return
}

// ------------------------------------------------------------------------------------------

