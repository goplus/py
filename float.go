package py

// #include <Python.h>
// static inline double floatCheck(PyObject *o) { return PyFloat_Check(o); }
import "C"
import "unsafe"

type Float struct {
	Base
	NumberProtocol
	o C.PyFloatObject
}

// FloatType is the Type object that represents the Float type.
var FloatType = (*Type)(unsafe.Pointer(&C.PyFloat_Type))

func newFloat(obj *C.PyObject) *Float {
	return (*Float)(unsafe.Pointer(obj))
}

func NewFloat(i float64) *Float {
	return newFloat(C.PyFloat_FromDouble(C.double(i)))
}

func AsFloat(o *Base) (v *Float, ok bool) {
	if ok = C.floatCheck(o.c()) != 0; ok {
		v = newFloat(o.c())
	}
	return
}

func NewFloatFromString(s string) *Float {
	cs := NewString(s)
	return newFloat(C.PyFloat_FromString((*C.PyObject)(unsafe.Pointer(cs.Obj())), nil))
}

func (f *Float) Float() float64 {
	return float64(C.PyFloat_AsDouble(f.c()))
}
