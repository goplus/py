package pyutil

import (
	"syscall"
	"github.com/qiniu/py"
	"github.com/qiniu/errors"
)

// ------------------------------------------------------------------------------------------

func PackEx(cfg *Config, args ...interface{}) (ret *py.Tuple, err error) {

	args1 := py.NewTuple(len(args))

	for i, arg := range args {
		v1, ok1 := NewVarEx(arg, cfg)
		if !ok1 {
			args1.Decref()
			err = errors.Info(syscall.EINVAL, "pyutil.Pack", i+1, arg).Detail(err)
			return
		}
		args1.SetItem(i, v1)
	}
	return args1, nil
}

func Pack(args ...interface{}) (ret *py.Tuple, err error) {

	return PackEx(DefaultConfig, args...)
}

// ------------------------------------------------------------------------------------------

func CallEx(cfg *Config, fn *py.Base, args ...interface{}) (ret *py.Base, err error) {

	args1, err := PackEx(cfg, args...)
	if err != nil {
		err = errors.Info(syscall.EINVAL, "pyutil.Call").Detail(err)
		return
	}
	defer args1.Decref()

	return fn.CallObject(args1)
}

func Call(fn *py.Base, args ...interface{}) (*py.Base, error) {

	return CallEx(DefaultConfig, fn, args...)
}

// ------------------------------------------------------------------------------------------

func CallMethodEx(cfg *Config, self *py.Base, method string, args ...interface{}) (ret *py.Base, err error) {

	args1, err := PackEx(cfg, args...)
	if err != nil {
		err = errors.Info(syscall.EINVAL, "pyutil.Call").Detail(err)
		return
	}
	defer args1.Decref()

	return self.CallMethodObject(method, args1)
}

func CallMethod(self *py.Base, method string, args ...interface{}) (ret *py.Base, err error) {

	return CallMethodEx(DefaultConfig, self, method, args...)
}

// ------------------------------------------------------------------------------------------

func NewInstanceEx(cfg *Config, typ *py.Class, args ...interface{}) (ret *py.Base, err error) {

	args1, err := PackEx(cfg, args...)
	if err != nil {
		err = errors.Info(syscall.EINVAL, "pyutil.NewInstance").Detail(err)
		return
	}
	defer args1.Decref()

	return typ.New(args1, nil)
}

func NewInstance(typ *py.Class, args ...interface{}) (ret *py.Base, err error) {

	return NewInstanceEx(DefaultConfig, typ, args...)
}

// ------------------------------------------------------------------------------------------

func NewEx(cfg *Config, mod *py.Base, clsname string, args ...interface{}) (ret *py.Base, err error) {

	o, err := mod.GetAttrString(clsname)
	if err != nil {
		err = errors.Info(err, "pyutil.New", clsname).Detail(err)
		return
	}
	defer o.Decref()

	ty, ok := py.AsClass(o)
	if !ok {
		err = errors.Info(syscall.EINVAL, "pyutil.New", o.String(), "is not a class")
		return
	}

	return NewInstanceEx(cfg, ty, args...)
}

func New(mod *py.Base, clsname string, args ...interface{}) (ret *py.Base, err error) {

	return NewEx(DefaultConfig, mod, clsname, args...)
}

// ------------------------------------------------------------------------------------------

