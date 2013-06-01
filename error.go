package py

// #include <Python.h>
// static inline void incref(PyObject *obj) { Py_INCREF(obj); }
// static inline void decref(PyObject *obj) { Py_DECREF(obj); }
// static inline void xdecref(PyObject *obj) { Py_XDECREF(obj); }
import "C"
import "fmt"
import "syscall"
import "strings"
import "runtime"
import "github.com/qiniu/errors"

// Error represents a Python exception as a Go struct that implements the
// error interface.  It allows Go code to handle Python exceptions in an
// idiomatic Go fashion.
type Error struct {
	Kind  *Base
	Value *Base
	tb    *C.PyObject
}

func newError(kind, val *Base, tb *C.PyObject) *Error {
	e := &Error{kind, val, tb}
	runtime.SetFinalizer(e, (*Error).release)
	return e
}

func (e *Error) release() error {
	if e.Kind != nil {
		e.Kind.Decref()
		e.Value.Decref()
		e.Kind = nil
		e.Value = nil
		if e.tb != nil {
			C.decref(e.tb)
			e.tb = nil
		}
	}
	return nil
}

// Error() returns a string representation of the Python exception represented
// by the Error e.  This is the same as the final line of the Python output from
// an uncaught exception.
func (e *Error) Error() string {
	kind := e.Kind.String()
	if strings.HasPrefix(kind, "<type 'exceptions.") {
		kind = kind[18 : len(kind)-2]
	}
	return kind + ": " + e.Value.String()
}

/*
// Matches returns true if e.Kind matches the exception in exc.  If exc is a
// Class, then true is returned if e.Kind is an instance.  If exc is a Tuple,
// then all elements (and recursively for sub elements) are searched for a
// match.
func (e *Error) Matches(exc Object) bool {
	return C.PyErr_GivenExceptionMatches(c(e.Kind), c(exc)) != 0
}
*/

// Normalize adjusts e.Kind/e.Value in the case that the values aren't
// normalized to start with.  It's possible that an Error returned from Python
// might have e.Kind be a Class, with e.Value not being an instance of that
// class, Normalize will fix this.  The separate normalization is implemented in
// Python to improve performance.
func (e *Error) Normalize() {
	exc := e.Kind.c()
	val := e.Value.c()
	tb := e.tb
	C.PyErr_NormalizeException(&exc, &val, &tb)
	e.Kind = newObject(exc)
	e.Value = newObject(val)
	e.tb = tb
}

// NewErrorV returns a new Error of the specified kind, and with the given
// value.
func NewErrorV(kind *Base, value *Base) *Error {
	kind.Incref()
	value.Incref()
	return newError(kind, value, nil)
}

// NewError returns a new Error of the specified kind, and with the value
// being a new String containing the string created the given format and args.
func NewError(kind *Base, format string, args ...interface{}) *Error {
	msg := fmt.Sprintf(format, args...)
	kind.Incref()
	val := NewString(msg)
	return newError(kind, &val.Base, nil)
}

func Raise(err error) {

	var val *C.PyObject
	var exc = C.PyExc_Exception

	e, ok := err.(*Error)
	if ok {
		exc = e.Kind.c()
		val = e.Value.c()
	} else {
		v := NewString(errors.Detail(err))
		val = v.c()
		defer C.decref(val)
	}
	C.PyErr_SetObject(exc, val)
}

func GetException() error {
	return exception()
}

func exceptionRaised() bool {
	return C.PyErr_Occurred() != nil
}

func exception() error {
	if C.PyErr_Occurred() == nil {
		return syscall.EFAULT
	}

	var t, v, tb *C.PyObject
	C.PyErr_Fetch(&t, &v, &tb)

	return newError(newObject(t), newObject(v), tb)
}

func ssize_t2Int64Err(s C.Py_ssize_t) (int64, error) {
	if s < 0 {
		return 0, exception()
	}
	return int64(s), nil
}

func int2BoolErr(i C.int) (bool, error) {
	if i < 0 {
		return false, exception()
	}
	return i > 0, nil
}

func int2Err(i C.int) error {
	if i < 0 {
		return exception()
	}
	return nil
}

func obj2ObjErr(obj *C.PyObject) (*Base, error) {
	if obj == nil {
		return nil, exception()
	}
	return newObject(obj), nil
}
