package py

/*
#include <Python.h>
#include "gofunction.h"

//static inline void decref(PyObject *obj) { Py_DECREF(obj); }
*/
import "C"
import "unsafe"
import "reflect"

type Closure struct { // closure = self.method
	Self      reflect.Value
	Method    reflect.Value
	methodDef C.PyMethodDef
}

func (closure *Closure) NewFunction(name string, nin int, doc string) *Base {

	d := &closure.methodDef
	d.ml_name = C.CString(name)
	defer C.free(unsafe.Pointer(d.ml_name))

	if C.setMethod(d, C.int(nin)) != 0 {
		panic("Invalid arguments: nin")
	}
	if doc != "" {
		d.ml_doc = C.CString(doc)
		defer C.free(unsafe.Pointer(d.ml_doc))
	}

	ctx := uintptr(unsafe.Pointer(closure))
	self := C.PyLong_FromLongLong(C.longlong(ctx))
	defer C.decref(self)

	f := C.PyCFunction_NewEx(d, self, nil)
	return (*Base)(unsafe.Pointer(f))
}

//export goClassCallMethodArgs
func goClassCallMethodArgs(obj, args unsafe.Pointer) unsafe.Pointer {

	// Unpack context and self pointer from obj
	t := (*C.PyObject)(obj)
	closure := (*Closure)(unsafe.Pointer(uintptr(C.PyLong_AsLongLong(t))))

	// Get args ready to use, by turning it into a pointer of the appropriate
	// type
	a := (*Tuple)(args)

	in := []reflect.Value{closure.Self, reflect.ValueOf(a)}
	out := closure.Method.Call(in)

	err := out[1].Interface()
	if err != nil {
		Raise(err.(error))
		return nil
	}

	ret := out[0].Interface().(*Base)
	return unsafe.Pointer(ret)
}

//export goClassCallMethodKwds
func goClassCallMethodKwds(obj, args, kwds unsafe.Pointer) unsafe.Pointer {

	// Unpack context and self pointer from obj
	t := (*C.PyObject)(obj)
	closure := (*Closure)(unsafe.Pointer(uintptr(C.PyLong_AsLongLong(t))))

	// Get args and kwds ready to use, by turning them into pointers of the
	// appropriate type
	a := (*Tuple)(args)
	k := (*Dict)(kwds)

	in := []reflect.Value{closure.Self, reflect.ValueOf(a), reflect.ValueOf(k)}
	out := closure.Method.Call(in)

	err := out[1].Interface()
	if err != nil {
		Raise(err.(error))
		return nil
	}

	ret := out[0].Interface().(*Base)
	return unsafe.Pointer(ret)
}
