package py

// #include <Python.h>
// static inline int typeCheck(PyObject *o) { return PyType_Check(o); }
// static inline int typeCheckE(PyObject *o) { return PyType_CheckExact(o); }
// static inline PyObject* typeAlloc(PyObject *t, Py_ssize_t n) {
//    return ((PyTypeObject*)t)->tp_alloc((PyTypeObject *)t, n);
// }
// static inline int typeInit(PyObject *t, PyObject *o, PyObject *a, PyObject *k) {
//    return ((PyTypeObject*)t)->tp_init(o, a, k);
// }
// static inline PyObject* typeNew(PyObject *t, PyObject *a, PyObject *k) {
//    return ((PyTypeObject*)t)->tp_new((PyTypeObject*)t, a, k);
// }
import "C"
import "unsafe"

type Type struct {
	Base
	o C.PyTypeObject
}

// TypeType is the Type object that represents the Type type.
var TypeType = (*Type)(unsafe.Pointer(&C.PyType_Type))

func newType(obj *C.PyObject) *Type {
	return (*Type)(unsafe.Pointer(obj))
}

func AsType(o *Base) (v *Type, ok bool) {
	if ok = C.typeCheck(o.c()) != 0; ok {
		v = newType(o.c())
	}
	return
}

func (t *Type) NewNoArgs() (ret *Base, err error) {
	args := NewTuple(0)
	defer args.Decref()
	return t.New(args, nil)
}

func (t *Type) New(args *Tuple, kw *Dict) (ret *Base, err error) {
	ret1 := C.typeNew(t.c(), args.c(), kw.c())
	return obj2ObjErr(ret1)
}

func (t *Type) NewObjArgs(args ...*Base) (ret *Base, err error) {
	args1 := PackTuple(args...)
	defer args1.Decref()
	return t.New(args1, nil)
}

func (t *Type) Alloc(n int64) (*Base, error) {
	ret := C.typeAlloc(t.c(), C.Py_ssize_t(n))
	return obj2ObjErr(ret)
}

func (t *Type) Init(obj *Base, args *Tuple, kw *Dict) error {
	ret := C.typeInit(t.c(), obj.c(), args.c(), kw.c())
	if ret < 0 {
		return exception()
	}
	return nil
}

// CheckExact returns true when "t" is an actual Type object, and not some form
// of subclass.
func (t *Type) CheckExact() bool {
	return C.typeCheckE(t.c()) == 1
}

// PyType_ClearCache : TODO - ???

// Modified should be called after the attributes or base class of a Type have
// been changed.
func (t *Type) Modified() {
	C.PyType_Modified(&t.o)
}

// HasFeature returns true when "t" has the feature in question.
func (t *Type) HasFeature(feature uint32) bool {
	return (t.o.tp_flags & C.long(feature)) != 0
}

// IsGc returns true if the type "t" supports Cyclic Garbage Collection.
func (t *Type) IsGc() bool {
	return t.HasFeature(C.Py_TPFLAGS_HAVE_GC)
}

// IsSubtype returns true if "t" is a subclass of "t2".
func (t *Type) IsSubtype(t2 *Type) bool {
	return C.PyType_IsSubtype(&t.o, &t2.o) == 1
}

// PyType_GenericAlloc : This is an internal function, which we should not need
// to expose.

// PyType_GenericNew : Another internal function we don't need to expose.

// PyType_Ready : This function is wrapped (along with a lot of other
// functionality) in the Create method of the Class stuct.
