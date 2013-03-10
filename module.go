package py

// #include <Python.h>
// static inline int moduleCheck(PyObject *o) { return PyModule_Check(o); }
// static inline int moduleCheckE(PyObject *o) { return PyModule_CheckExact(o); }
// static inline void decref(PyObject *obj) { Py_DECREF(obj); }
import "C"
import "unsafe"

type Module struct {
	Base
	o C.PyObject
}

// ModuleType is the Type object that represents the Module type.
var ModuleType = (*Type)(unsafe.Pointer(&C.PyModule_Type))

func newModule(obj *C.PyObject) *Module {
	return (*Module)(unsafe.Pointer(obj))
}

func Import(name string) (*Module, error) {
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))

	pyName := C.PyString_FromString(s)
	defer C.decref(pyName)

	obj := C.PyImport_Import(pyName)
	if obj == nil {
		return nil, exception()
	}

	return newModule(obj), nil
}

func ExecCodeModule(name string, code *Base) (*Module, error) {
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))
	ret := C.PyImport_ExecCodeModule(s, code.c())
	if ret == nil {
		return nil, exception()
	}
	return newModule(ret), nil
}

func NewModule(name string) (*Module, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	ret := C.PyModule_New(cname)
	if ret == nil {
		return nil, exception()
	}

	return newModule(ret), nil
}

func AsModule(o *Base) (v *Module, ok bool) {
	if ok = C.moduleCheck(o.c()) != 0; ok {
		v = newModule(o.c())
	}
	return
}

func (mod *Module) CheckExact() bool {
	return C.moduleCheckE(mod.c()) != 0
}

// Return value: Borrowed reference.
func (mod *Module) Dict() *Dict {
	ret := C.PyModule_GetDict(mod.c())
	return newDict(ret)
}

// Return module‘s __name__ value. If the module does not provide one, or if it is not a string, 
// SystemError is raised and NULL is returned.
func (mod *Module) Name() (string, error) {
	ret := C.PyModule_GetName(mod.c())
	if ret == nil {
		return "", exception()
	}
	return C.GoString(ret), nil
}

func (mod *Module) Filename() (string, error) {
	ret := C.PyModule_GetFilename(mod.c())
	if ret == nil {
		return "", exception()
	}
	return C.GoString(ret), nil
}

func (mod *Module) AddObject(name string, obj *Base) error {
	if obj == nil {
		return AssertionError.Err("ValueError: obj == nil!")
	}

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	ret := C.PyModule_AddObject(mod.c(), cname, obj.c())
	if ret < 0 {
		return exception()
	}

	return nil
}

func (mod *Module) AddIntConstant(name string, value int) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	ret := C.PyModule_AddIntConstant(mod.c(), cname, C.long(value))
	if ret < 0 {
		return exception()
	}

	return nil
}

func (mod *Module) AddStringConstant(name, value string) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))

	ret := C.PyModule_AddStringConstant(mod.c(), cname, cvalue)
	if ret < 0 {
		return exception()
	}

	return nil
}
