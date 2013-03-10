package py

// #include <Python.h>
// static inline long longCheck(PyObject *o) { return PyLong_Check(o); }
import "C"
import "unsafe"

type Long struct {
	Base
	NumberProtocol
	o C.PyLongObject
}

// LongType is the Type object that represents the Long type.
var LongType = (*Type)(unsafe.Pointer(&C.PyLong_Type))

func newLong(obj *C.PyObject) *Long {
	return (*Long)(unsafe.Pointer(obj))
}

func NewLong(i int64) *Long {
	return newLong(C.PyLong_FromLongLong(C.longlong(i)))
}

func AsLong(o *Base) (v *Long, ok bool) {
	if ok = C.longCheck(o.c()) != 0; ok {
		v = newLong(o.c())
	}
	return
}

func (l *Long) Long() int64 {
	return int64(C.PyLong_AsLongLong(l.c()))
}

