package py

// #include <Python.h>
// static inline int tupleCheckE(PyObject *o) { return PyTuple_CheckExact(o); }
// static inline int tupleCheck(PyObject *o) { return PyTuple_Check(o); }
// static inline size_t tupleItemSize() { return sizeof(PyObject *); }
import "C"
import "unsafe"

type Tuple struct {
	Base
	o C.PyTupleObject
}

// TupleType is the Type object that represents the Tuple type.
var TupleType = (*Type)(unsafe.Pointer(&C.PyTuple_Type))

func newTuple(obj *C.PyObject) *Tuple {
	return (*Tuple)(unsafe.Pointer(obj))
}

// NewTuple returns a new *Tuple of the specified size.  However the entries are
// all set to NULL, so the tuple should not be shared, especially with Python
// code, until the entries have all been set.
//
// Return value: New Reference.
func NewTuple(size int) *Tuple {
	ret := C.PyTuple_New(C.Py_ssize_t(size))
	return newTuple(ret)
}

func AsTuple(o *Base) (v *Tuple, ok bool) {
	if ok = C.tupleCheck(o.c()) != 0; ok {
		v = newTuple(o.c())
	}
	return
}

// PackTuple returns a new *Tuple which contains the arguments.  This tuple is
// ready to use.
//
// Return value: New Reference.
func PackTuple(items ...*Base) *Tuple {
	ret := C.PyTuple_New(C.Py_ssize_t(len(items)))

	// Since the ob_item array has a size of 1, Go won't let us index more than
	// a single entry, and if we try and use our own local type definition with
	// a flexible array member then cgo converts it to [0]byte which is even
	// less useful.  So, we resort to pointer manipulation - which is
	// unfortunate, as it's messy in Go.

	// base is a pointer to the first item in the array of PyObject pointers.
	// step is the size of a PyObject * (i.e. the number of bytes we need to add
	// to get to the next item).
	base := unsafe.Pointer(&(*C.PyTupleObject)(unsafe.Pointer(ret)).ob_item[0])
	step := uintptr(C.tupleItemSize())

	for _, item := range items {
		item.Incref()
		*(**C.PyObject)(base) = item.c()

		// Move base to point to the next item, by incrementing by step bytes
		base = unsafe.Pointer(uintptr(base) + step)
	}
	return newTuple(ret)
}

func (t *Tuple) CheckExact() bool {
	ret := C.tupleCheckE(t.c())
	if int(ret) != 0 {
		return true
	}
	return false
}

func (t *Tuple) Size() int {
	ret := C.PyTuple_Size(t.c())
	if ret < 0 {
		panic(exception())
	}
	return int(ret)
}

// Return the object at position pos in the tuple pointed to by p. If pos is out of bounds, 
// return NULL and sets an IndexError exception.
//
// Return value: Borrowed reference.
func (t *Tuple) GetItem(pos int) (*Base, error) {
	ret := C.PyTuple_GetItem(t.c(), C.Py_ssize_t(pos))
	return obj2ObjErr(ret)
}

func (t *Tuple) GetSlice(low, high int) (*Tuple, error) {
	ret := C.PyTuple_GetSlice(t.c(), C.Py_ssize_t(low), C.Py_ssize_t(high))
	if ret == nil {
		return nil, exception()
	}
	return newTuple(ret), nil
}

// Insert a reference to object o at position pos of the tuple pointed to by p. Return 0 on success.
// Note This function “steals” a reference to o.
func (t *Tuple) SetItem(pos int, obj *Base) error {
	ret := C.PyTuple_SetItem(t.c(), C.Py_ssize_t(pos), obj.c())
	return int2Err(ret)
}

// _PyTuple_Resize

// PyTuple_ClearFreeList()

func (t *Tuple) Slice() []*Base {
	l := t.Size()
	s := make([]*Base, l)
	for i := 0; i < l; i++ {
		o, err := t.GetItem(i)
		if err != nil {
			panic(err)
		}
		s[i] = o
	}
	return s
}

