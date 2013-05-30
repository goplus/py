package py

/*
#include <Python.h>

static inline PyObject* incNone() { Py_RETURN_NONE; }
*/
import "C"

// ------------------------------------------------------------------------------------------

func IncNone() *Base {
	return newObject(C.incNone())
}

// ------------------------------------------------------------------------------------------
