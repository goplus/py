package py

// #include <Python.h>
// static inline int stringCheck(PyObject *o) { return PyString_Check(o); }
import "C"
import "unsafe"

type String struct {
	Base
	o C.PyStringObject
}

// StringType is the Type object that represents the String type.
var StringType = (*Type)(unsafe.Pointer(&C.PyString_Type))

func newString(obj *C.PyObject) *String {
	return (*String)(unsafe.Pointer(obj))
}

func NewString(s string) *String {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	ret := C.PyString_FromString(cs)
	return newString(ret)
}

func AsString(o *Base) (v *String, ok bool) {
	if ok = C.stringCheck(o.c()) != 0; ok {
		v = newString(o.c())
	}
	return
}

func (s *String) String() string {
	if s == nil {
		return "<nil>"
	}
	ret := C.PyString_AsString(s.c())
	return C.GoString(ret)
}

func (s *String) Format(args *Tuple) (*String, error) {
	ret := C.PyString_Format(s.c(), args.c())
	if ret == nil {
		return nil, exception()
	}
	return newString(ret), nil
}
