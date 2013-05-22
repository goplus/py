package py

// #include <Python.h>
// static inline long floatCheck(PyObject *o) { return PyFloat_Check(o); }
import "C"
import "unsafe"

type Float struct {
	Base
	NumberProtocol
	o C.PyFloatObject
}

// LongType is the Type object that represents the Long type.
var FloatType = (*Type)(unsafe.Pointer(&C.PyFloat_Type))

func newFloat(obj *C.PyObject) *Float {
	return (*Float)(unsafe.Pointer(obj))
}

func NewFloat(i float64) *Float {
	return newFloat(C.PyLong_FromDouble(C.double(i)))
}

func AsFloat(o *Base) (v *Float, ok bool) {
	if ok = C.floatCheck(o.c()) != 0; ok {
		v = newFloat(o.c())
	}
	return
}

func (f *Float) Float() float64 {
	return float64(C.PyLong_AsDouble(f.c()))
}

