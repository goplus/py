package py

// #include <Python.h>
// static inline int dictCheck(PyObject *o) { return PyDict_Check(o); }
// static inline int dictCheckE(PyObject *o) { return PyDict_CheckExact(o); }
import "C"
import "unsafe"

// *Dict represents a Python dictionary.  In addition to satisfying the Object
// interface, Dict pointers also have a number of methods defined - representing
// the PyDict_XXX functions from the Python C API.
type Dict struct {
	Base
	o C.PyDictObject
}

// DictType is the Type object that represents the Dict type.
var DictType = (*Type)(unsafe.Pointer(&C.PyDict_Type))

func newDict(obj *C.PyObject) *Dict {
	return (*Dict)(unsafe.Pointer(obj))
}

// NewDict creates a new empty dictionary.
//
// Return value: New Reference.
func NewDict() *Dict {
	ret := C.PyDict_New()
	return newDict(ret)
}

func NewDictProxy(obj *Base) *Dict {
	ret := C.PyDictProxy_New(obj.c())
	return newDict(ret)
}

func AsDict(o *Base) (v *Dict, ok bool) {
	if ok = C.dictCheck(o.c()) != 0; ok {
		v = newDict(o.c())
	}
	return
}

// CheckExact returns true if d is an actual dictionary object, and not an
// instance of a sub type.
func (d *Dict) CheckExact() bool {
	ret := C.dictCheckE(d.c())
	if int(ret) != 0 {
		return true
	}
	return false
}

// Clear empties the dictionary d of all key-value pairs.
func (d *Dict) Clear() {
	C.PyDict_Clear(d.c())
}

// Contains Returns true if the dictionary contains the given key.  This is
// equivalent to the Python expression "key in d".
func (d *Dict) Contains(key *Base) (bool, error) {
	ret := C.PyDict_Contains(d.c(), key.c())
	return int2BoolErr(ret)
}

// Copy returns a new dictionary that contains the same key-values pairs as d.
//
// Return value: New Reference.
func (d *Dict) Copy() (*Base, error) {
	ret := C.PyDict_Copy(d.c())
	return obj2ObjErr(ret)
}

// SetItem inserts "val" into dictionary d with the key "key".  If "key" is not
// hashable, then a TypeError will be returned.
func (d *Dict) SetItem(key, val *Base) error {
	ret := C.PyDict_SetItem(d.c(), key.c(), val.c())
	return int2Err(ret)
}

// SetItemString inserts "val" into dictionary d with the key "key" (or rather,
// with a *String with the value of "key" will be used as the key).  If "key" is
// not hashable, then a TypeError will be returned.
func (d *Dict) SetItemString(key string, val *Base) error {
	s := C.CString(key)
	defer C.free(unsafe.Pointer(s))
	ret := C.PyDict_SetItemString(d.c(), s, val.c())
	return int2Err(ret)
}

// DelItem removes the entry with the key of "key" from the dictionary d.  If
// "key" is not hashable, a TypeError is returned.
func (d *Dict) DelItem(key *Base) error {
	ret := C.PyDict_DelItem(d.c(), key.c())
	return int2Err(ret)
}

// DelItem removes the entry with the key of "key" (or rather, with a *String
// with the value of "key" as the key) from the dictionary d.
func (d *Dict) DelItemString(key string) error {
	s := C.CString(key)
	defer C.free(unsafe.Pointer(s))
	ret := C.PyDict_DelItemString(d.c(), s)
	return int2Err(ret)
}

// GetItem returns the Object from dictionary d which has the key "key".  If
// there is no such Object, then nil is returned (without an error).
//
// Return value: Borrowed Reference.
func (d *Dict) GetItem(key *Base) *Base {
	ret := C.PyDict_GetItem(d.c(), key.c())
	return newObject(ret)
}

// GetItemString returns the Object from dictionary d which has the key "key"
// (or rather, which has a *String with the value of "key" as the key).  If
// there is no such Object, then nil is returned (without an error).
//
// Return value: Borrowed Reference.
func (d *Dict) GetItemString(key string) *Base {
	s := C.CString(key)
	defer C.free(unsafe.Pointer(s))
	ret := C.PyDict_GetItemString(d.c(), s)
	return newObject(ret)
}

/*
// Items returns a *List containing all the items from the dictionary d, as with
// the Python "d.items()".
//
// Return value: New Reference.
func (d *Dict) Items() (*List, error) {
	ret := C.PyDict_Items(d.c())
	return newList(ret), exception()
}

// Keys returns a *List containing all the keys from the dictionary d, as with
// the Python "d.keys()".
//
// Return value: New Reference.
func (d *Dict) Keys() (*List, error) {
	ret := C.PyDict_Keys(d.c())
	return newList(ret), exception()
}

// Values returns a *List containing all the values from the dictionary d, as
// with the Python "d.values()".
//
// Return value: New Reference.
func (d *Dict) Values() (*List, error) {
	ret := C.PyDict_Values(d.c())
	return newList(ret), exception()
}
*/

// Size returns the number of items in the dictionary d.  This is equivalent to
// the Python "len(d)".
func (d *Dict) Size() int {
	ret := C.PyDict_Size(d.c())
	if ret < 0 {
		panic(exception())
	}
	return int(ret)
}

// PyDict_Next

// Merge merges key values pairs from Object o (which may be a dictionary, or an
// object that supports "o.keys()" and "o[key]") into the dictionary d.  If
// override is true then a matching key in d will have it's value replaced by
// the one in o, else the value in d will be left.
func (d *Dict) Merge(o *Base, override bool) error {
	over := 0
	if override {
		over = 1
	}
	ret := C.PyDict_Merge(d.c(), o.c(), C.int(over))
	return int2Err(ret)
}

// Update replaces key values pairs in d with those from o.  It is equivalent to
// d.Merge(o, true) in Go, or "d.update(o)" in Python.
func (d *Dict) Update(o *Base) error {
	ret := C.PyDict_Update(d.c(), o.c())
	return int2Err(ret)
}

// MergeFromSeq2 merges key values pairs from the Object o (which must be an
// iterable object, where each item is an iterable of length 2 - the key value
// pairs).  If override is true then the last key value pair with the same key
// wins, otherwise the first instance does (where an instance already in d
// counts before any in o).
func (d *Dict) MergeFromSeq2(o *Base, override bool) error {
	over := 0
	if override {
		over = 1
	}
	ret := C.PyDict_MergeFromSeq2(d.c(), o.c(), C.int(over))
	return int2Err(ret)
}

// Map returns a Go map that contains the values from the Python dictionary,
// indexed by the keys.  The keys and values are the same as in the Python
// dictionary, but changes to the Go map are not propogated back to the Python
// dictionary.
//
// Note: the map holds borrowed references
func (d *Dict) Map() map[*Base]*Base {
	m := make(map[*Base]*Base, d.Size())
	var p C.Py_ssize_t
	var k *C.PyObject
	var v *C.PyObject
	for int(C.PyDict_Next(d.c(), &p, &k, &v)) != 0 {
		key := newObject(k)
		value := newObject(v)
		m[key] = value
	}
	return m
}

// MapString is similar to Map, except that the keys are first converted to
// strings.  If the keys are not all Python strings, then an error is returned.
//
// Note: the map holds borrowed references
func (d *Dict) MapString() (map[string]*Base, error) {
	m := make(map[string]*Base, d.Size())
	var p DictIter
	var k, v *Base
	for d.Next(&p, &k, &v) {
		s, ok := AsString(k)
		if !ok {
			return nil, TypeError.Err("%v is not a string", k)
		}
		m[s.String()] = v
	}
	return m, nil
}

type DictIter C.Py_ssize_t

// Iterate over all key-value pairs in the dictionary d.
// The Py_ssize_t referred to by ppos must be initialized to 0 prior to the first call to this function
// to start the iteration; the function returns true for each pair in the dictionary, and false once all
// pairs have been reported. The parameters pkey and pvalue should either point to PyObject* variables
// that will be filled in with each key and value, respectively, or may be NULL. Any references returned
// through them are borrowed. ppos should not be altered during iteration. Its value represents offsets
// within the internal dictionary structure, and since the structure is sparse, the offsets are not consecutive.
func (d *Dict) Next(pos *DictIter, k, v **Base) bool {
	k1 := (**C.PyObject)(unsafe.Pointer(k))
	v1 := (**C.PyObject)(unsafe.Pointer(v))
	return C.PyDict_Next(d.c(), (*C.Py_ssize_t)(pos), k1, v1) != 0
}
