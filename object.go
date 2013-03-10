package py

// #include <Python.h>
// static inline void incref(PyObject *obj) { Py_INCREF(obj); }
// static inline void decref(PyObject *obj) { Py_DECREF(obj); }
import "C"
import "unsafe"

type Op int

const (
	LT = Op(C.Py_LT)
	LE = Op(C.Py_LE)
	EQ = Op(C.Py_EQ)
	NE = Op(C.Py_NE)
	GT = Op(C.Py_GT)
	GE = Op(C.Py_GE)
)

// Base is an 0-sized type that can be embedded as the first item in
// concrete types to provide the Object interface functions.
type Base struct {}

func newObject(obj *C.PyObject) *Base {
	return (*Base)(unsafe.Pointer(obj))
}

func (obj *Base) c() *C.PyObject {
	return (*C.PyObject)(unsafe.Pointer(obj))
}

func (obj *Base) Obj() *Base {
	return obj
}

// Init initialises obj.  It is equivalent to "obj.__init__(*args, **kw)" in
// Python.
func (obj *Base) Init(args *Tuple, kw *Dict) error {
	return obj.Type().Init(obj, args, kw)
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (obj *Base) Type() *Type {
	o := obj.c().ob_type
	return newType((*C.PyObject)(unsafe.Pointer(o)))
}

// Decref decrements obj's reference count, obj may not be nil.
func (obj *Base) Decref() {
	C.decref(obj.c())
}

// Incref increments obj's reference count, obj may not be nil.
func (obj *Base) Incref() {
	C.incref(obj.c())
}

// IsTrue returns true if the value of obj is considered to be True.  This is
// equivalent to "if obj:" in Python.
func (obj *Base) IsTrue() bool {
	ret := C.PyObject_IsTrue(obj.c())
	if ret < 0 {
		panic(exception())
	}
	return ret != 0
}

// Not returns true if the value of obj is considered to be False.  This is
// equivalent to "if not obj:" in Python.
func (obj *Base) Not() bool {
	ret := C.PyObject_Not(obj.c())
	if ret < 0 {
		panic(exception())
	}
	return ret != 0
}

// HasAttr returns true if "obj" has the attribute "name".  This is equivalent
// to the Python "hasattr(obj, name)".
func (obj *Base) HasAttr(name *Base) bool {
	ret := C.PyObject_HasAttr(obj.c(), name.c())
	if ret == 1 {
		return true
	}
	return false
}

// HasAttrString returns true if "obj" has the attribute "name".  This is
// equivalent to the Python "hasattr(obj, name)".
func (obj *Base) HasAttrString(name string) bool {
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))
	ret := C.PyObject_HasAttrString(obj.c(), s)
	if ret == 1 {
		return true
	}
	return false
}

// GetAttr returns the attribute of "obj" with the name "name".  This is
// equivalent to the Python "obj.name".
//
// Return value: New Reference.
func (obj *Base) GetAttr(name *Base) (*Base, error) {
	ret := C.PyObject_GetAttr(obj.c(), name.c())
	return obj2ObjErr(ret)
}

// Retrieve an attribute named attr_name from object o. Returns the attribute value
// on success, or NULL on failure. This is the equivalent to the Python "obj.name".
//
// Return value: New reference.
func (obj *Base) GetAttrString(name string) (*Base, error) {
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))
	ret := C.PyObject_GetAttrString(obj.c(), s)
	return obj2ObjErr(ret)
}

// PyObject_GenericGetAttr : This is an internal helper function - we shouldn't
// need to expose it ...

// SetAttr sets the attribute of "obj" with the name "name" to "value".  This is
// equivalent to the Python "obj.name = value".
func (obj *Base) SetAttr(name, value *Base) error {
	ret := C.PyObject_SetAttr(obj.c(), name.c(), value.c())
	return int2Err(ret)
}

// SetAttrString sets the attribute of "obj" with the name "name" to "value".
// This is equivalent to the Python "obj.name = value".
func (obj *Base) SetAttrString(name string, value *Base) error {
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))
	ret := C.PyObject_SetAttrString(obj.c(), s, value.c())
	return int2Err(ret)
}

// PyObject_GenericSetAttr : This is an internal helper function - we shouldn't
// need to expose it ...

// DelAttr deletes the attribute with the name "name" from "obj".  This is
// equivalent to the Python "del obj.name".
func (obj *Base) DelAttr(name *Base) error {
	ret := C.PyObject_SetAttr(obj.c(), name.c(), nil)
	return int2Err(ret)
}

// DelAttrString deletes the attribute with the name "name" from "obj".  This is
// equivalent to the Python "del obj.name".
func (obj *Base) DelAttrString(name string) error {
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))
	ret := C.PyObject_SetAttrString(obj.c(), s, nil)
	return int2Err(ret)
}

// RichCompare compares "obj" with "obj2" using the specified operation (LE, GE
// etc.), and returns the result.  The equivalent Python is "obj op obj2", where
// op is the corresponding Python operator for op.
//
// Return value: New Reference.
func (obj *Base) RichCompare(obj2 *Base, op Op) (*Base, error) {
	ret := C.PyObject_RichCompare(obj.c(), obj2.c(), C.int(op))
	return obj2ObjErr(ret)
}

// RichCompare compares "obj" with "obj2" using the specified operation (LE, GE
// etc.), and returns true or false.  The equivalent Python is "obj op obj2",
// where op is the corresponding Python operator for op.
func (obj *Base) RichCompareBool(obj2 *Base, op Op) (bool, error) {
	ret := C.PyObject_RichCompareBool(obj.c(), obj2.c(), C.int(op))
	return int2BoolErr(ret)
}

// PyObject_Cmp : Thanks to multiple return values, we don't need this function
// to be available in Go.

// Compare returns the result of comparing "obj" and "obj2".  This is equivalent
// to the Python "cmp(obj, obj2)".
func (obj *Base) Compare(obj2 *Base) (int, error) {
	ret := C.PyObject_Compare(obj.c(), obj2.c())
	return int(ret), exception()
}

// Repr returns a String representation of "obj".  This is equivalent to the
// Python "repr(obj)".
//
// Return value: New Reference.
func (obj *Base) Repr() (*Base, error) {
	ret := C.PyObject_Repr(obj.c())
	return obj2ObjErr(ret)
}

// Str returns a String representation of "obj".  This is equivalent to the
// Python "str(obj)".
//
// Return value: New Reference.
func (obj *Base) Str() (*Base, error) {
	ret := C.PyObject_Str(obj.c())
	return obj2ObjErr(ret)
}

func (obj *Base) String() string {
	if v, ok := AsString(obj); ok {
		return v.String()
	}
	ret := C.PyObject_Str(obj.c())
	if ret == nil {
		return "<nil>"
	}
	defer C.decref(ret)
	return ((*String)(unsafe.Pointer(ret))).String()
}

// Bytes returns a Bytes representation of "obj".  This is equivalent to the
// Python "bytes(obj)".  In Python 2.x this method is identical to Str().
//
// Return value: New Reference.
func (obj *Base) Bytes() (*Base, error) {
	ret := C.PyObject_Bytes(obj.c())
	return obj2ObjErr(ret)
}

// PyObject_Unicode : TODO

// IsInstance returns true if "obj" is an instance of "cls", false otherwise.
// If "cls" is a Type instead of a class, then true will be return if "obj" is
// of that type.  If "cls" is a Tuple then true will be returned if "obj" is an
// instance of any of the Objects in the tuple.  This is equivalent to the
// Python "isinstance(obj, cls)".
func (obj *Base) IsInstance(cls *Base) (bool, error) {
	ret := C.PyObject_IsInstance(obj.c(), cls.c())
	return int2BoolErr(ret)
}

// IsSubclass retuns true if "obj" is a Subclass of "cls", false otherwise.  If
// "cls" is a Tuple, then true is returned if "obj" is a Subclass of any member
// of "cls".  This is equivalent to the Python "issubclass(obj, cls)".
func (obj *Base) IsSubclass(cls *Base) (bool, error) {
	ret := C.PyObject_IsSubclass(obj.c(), cls.c())
	return int2BoolErr(ret)
}

// Call calls obj with the given args and kwds.  kwds may be nil, args may not
// (an empty Tuple must be used if no arguments are wanted).  Returns the result
// of the call, or an Error on failure.  This is equivalent to
// "obj(*args, **kwds)" in Python.
//
// Return value: New Reference.
func (obj *Base) Call(args *Tuple, kwds *Dict) (*Base, error) {
	ret := C.PyObject_Call(obj.c(), args.c(), kwds.c())
	return obj2ObjErr(ret)
}

// CallObject calls obj with the given args.  args may be nil.  Returns the
// result of the call, or an Error on failure.  This is equivalent to
// "obj(*args)" in Python.
//
// Return value: New Reference.
func (obj *Base) CallObject(args *Tuple) (*Base, error) {
	var a *C.PyObject = nil
	if args != nil {
		a = args.c()
	}
	ret := C.PyObject_CallObject(obj.c(), a)
	return obj2ObjErr(ret)
}

func (obj *Base) CallMethodObject(name string, args *Tuple) (*Base, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	f := C.PyObject_GetAttrString(obj.c(), cname)
	if f == nil {
		return nil, AttributeError.Err(name)
	}
	defer C.decref(f)

	if C.PyCallable_Check(f) == 0 {
		return nil, TypeError.Err("attribute of type '%s' is not callable", name)
	}

	ret := C.PyObject_CallObject(f, args.c())
	return obj2ObjErr(ret)
}

func (obj *Base) CallObjArgs(args ...*Base) (*Base, error) {
	args1 := PackTuple(args...)
	defer args1.Decref()
	return obj.CallObject(args1)
}

func (obj *Base) CallMethodObjArgs(name string, args ...*Base) (*Base, error) {
	args1 := PackTuple(args...)
	defer args1.Decref()
	return obj.CallMethodObject(name, args1)
}

// PyObject_Hash : TODO

// PyObject_HashNotImplement : This is an internal function, that we probably
// don't need to export.

// Length returns the length of the Object.  This is equivalent to the Python
// "len(obj)".
func (obj *Base) Length() (int64, error) {
	ret := C.PyObject_Length(obj.c())
	return int64(ret), exception()
}

// Size returns the length of the Object.  This is equivalent to the Python
// "len(obj)".
func (obj *Base) Size() (int64, error) {
	ret := C.PyObject_Size(obj.c())
	return int64(ret), exception()
}

// GetItem returns the element of "obj" corresponding to "key".  This is
// equivalent to the Python "obj[key]".
//
// Return value: New Reference.
func (obj *Base) GetItem(key *Base) (*Base, error) {
	ret := C.PyObject_GetItem(obj.c(), key.c())
	return obj2ObjErr(ret)
}

// SetItem sets the element of "obj" corresponding to "key" to "value".  This is
// equivalent to the Python "obj[key] = value".
func (obj *Base) SetItem(key, value *Base) error {
	ret := C.PyObject_SetItem(obj.c(), key.c(), value.c())
	return int2Err(ret)
}

// DelItem deletes the element from "obj" that corresponds to "key".  This is
// equivalent to the Python "del obj[key]".
func (obj *Base) DelItem(key *Base) error {
	ret := C.PyObject_DelItem(obj.c(), key.c())
	return int2Err(ret)
}

// PyObject_AsFileDescriptor : TODO

func (obj *Base) Dir() (*Base, error) {
	ret := C.PyObject_Dir(obj.c())
	return obj2ObjErr(ret)
}

// PyObject_GetIter : TODO

