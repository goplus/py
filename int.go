// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include <Python.h>
// static inline int intCheck(PyObject *o) { return PyInt_Check(o); }
import "C"
import "unsafe"

type Int struct {
	Base
	NumberProtocol
	o C.PyIntObject
}

// IntType is the Type object that represents the Int type.
var IntType = (*Type)(unsafe.Pointer(&C.PyInt_Type))

func newInt(obj *C.PyObject) *Int {
	return (*Int)(unsafe.Pointer(obj))
}

func NewInt(i int) *Int {
	return newInt(C.PyInt_FromLong(C.long(i)))
}

func NewInt64(i int64) *Int {
	return newInt(C.PyInt_FromSsize_t(C.Py_ssize_t(i)))
}

func AsInt(o *Base) (v *Int, ok bool) {
	if ok = C.intCheck(o.c()) != 0; ok {
		v = newInt(o.c())
	}
	return
}

func (i *Int) Int() int {
	return int(C.PyInt_AsLong(i.c()))
}

