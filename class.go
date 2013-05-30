package py

// #include <Python.h>
// static inline int classCheck(PyObject *o) { return PyClass_Check(o); }
import "C"
import "unsafe"

type Class struct {
	Base
	o C.PyClassObject
}

// ClassType is the Type object that represents the Class type.
var ClassType = (*Type)(unsafe.Pointer(&C.PyClass_Type))

func newClass(obj *C.PyObject) *Class {
	return (*Class)(unsafe.Pointer(obj))
}

func AsClass(o *Base) (v *Class, ok bool) {
	if ok = C.classCheck(o.c()) != 0; ok {
		v = newClass(o.c())
	}
	return
}

func (t *Class) NewNoArgs() (ret *Base, err error) {
	args := NewTuple(0)
	defer args.Decref()
	return t.New(args, nil)
}

// Return value: New reference.
// Create a new instance of a specific class. The parameters arg and kw are used as
// the positional and keyword parameters to the object’s constructor.
func (t *Class) New(args *Tuple, kw *Dict) (ret *Base, err error) {
	ret1 := C.PyInstance_New(t.c(), args.c(), kw.c())
	return obj2ObjErr(ret1)
}

func (t *Class) NewObjArgs(args ...*Base) (ret *Base, err error) {
	args1 := PackTuple(args...)
	defer args1.Decref()
	return t.New(args1, nil)
}

// Return true if klass is a subclass of base. Return false in all other cases.
func (t *Class) IsSubclass(base *Base) bool {
	return C.PyClass_IsSubclass(t.c(), base.c()) != 0
}
