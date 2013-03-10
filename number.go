// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include <Python.h>
import "C"
import "unsafe"

// NumberProtocol is a 0-sized type that can be embedded in concrete types after
// the AbstractObject to provide access to the suite of methods that Python
// calls the "Number Protocol".
type NumberProtocol struct{}

func cnp(n *NumberProtocol) *C.PyObject {
	return (*C.PyObject)(unsafe.Pointer(n))
}

// Add returns the result of adding n and obj.  The equivalent Python is "n +
// obj".
//
// Return value: New Reference.
func (n *NumberProtocol) Add(obj *Base) (*Base, error) {
	ret := C.PyNumber_Add(cnp(n), obj.c())
	return obj2ObjErr(ret)
}

// Subtract returns the result of subtracting obj from n.  The equivalent Python
// is "n - obj".
//
// Return value: New Reference.
func (n *NumberProtocol) Subtract(obj *Base) (*Base, error) {
	ret := C.PyNumber_Subtract(cnp(n), obj.c())
	return obj2ObjErr(ret)
}

// Multiply returns the result of multiplying n by obj.  The equivalent Python
// is "n * obj".
//
// Return value: New Reference.
func (n *NumberProtocol) Multiply(obj *Base) (*Base, error) {
	ret := C.PyNumber_Multiply(cnp(n), obj.c())
	return obj2ObjErr(ret)
}

// Divide returns the result of dividing n by obj.  The equivalent Python is "n
// / obj".
//
// Return value: New Reference.
func (n *NumberProtocol) Divide(obj *Base) (*Base, error) {
	ret := C.PyNumber_Divide(cnp(n), obj.c())
	return obj2ObjErr(ret)
}

// FloorDivide returns the floor of dividing n obj obj.
//
// Return value: New Reference.
func (n *NumberProtocol) FloorDivide(obj *Base) (*Base, error) {
	ret := C.PyNumber_FloorDivide(cnp(n), obj.c())
	return obj2ObjErr(ret)
}

// TrueDivide returns the ... TODO
//
// Return value: New Reference.
func (n *NumberProtocol) TrueDivide(obj *Base) (*Base, error) {
	ret := C.PyNumber_TrueDivide(cnp(n), obj.c())
	return obj2ObjErr(ret)
}

// Remainder returns the remainder of dividing n by obj.  The equivalent Python
// is "n % obj".
//
// Return value: New Reference.
func (n *NumberProtocol) Remainder(obj *Base) (*Base, error) {
	ret := C.PyNumber_Remainder(cnp(n), obj.c())
	return obj2ObjErr(ret)
}

// Divmod returns the result of the Python "divmod(n, obj)".
//
// Return value: New Reference.
func (n *NumberProtocol) Divmod(obj *Base) (*Base, error) {
	ret := C.PyNumber_Divmod(cnp(n), obj.c())
	return obj2ObjErr(ret)
}

// Power returns the result of the Python "pow(n, obj1, obj2)".
//
// Return value: New Reference.
func (n *NumberProtocol) Power(obj1, obj2 *Base) (*Base, error) {
	ret := C.PyNumber_Power(cnp(n), obj1.c(), obj2.c())
	return obj2ObjErr(ret)
}

// Negative returns the negation of n.  The equivalent Python is "-n".
//
// Return value: New Reference.
func (n *NumberProtocol) Negative() (*Base, error) {
	ret := C.PyNumber_Negative(cnp(n))
	return obj2ObjErr(ret)
}

// Positive returns the positive of n.  The equivalent Python is "+n".
//
// Return value: New Reference.
func (n *NumberProtocol) Positive() (*Base, error) {
	ret := C.PyNumber_Positive(cnp(n))
	return obj2ObjErr(ret)
}

// Absolute returns the absolute value of n.  The equivalent Python is "abs(n)".
//
// Return value: New Reference.
func (n *NumberProtocol) Absolute() (*Base, error) {
	ret := C.PyNumber_Absolute(cnp(n))
	return obj2ObjErr(ret)
}

// Invert returns the bitwise negation of n.  The equivalent Python is "-n".
//
// Return value: New Reference.
func (n *NumberProtocol) Invert() (*Base, error) {
	ret := C.PyNumber_Invert(cnp(n))
	return obj2ObjErr(ret)
}

// Lshift returns the result of left shifting n by obj.  The equivalent Python
// is "n << obj".
//
// Return value: New Reference.
func (n *NumberProtocol) Lshift(obj *Base) (*Base, error) {
	ret := C.PyNumber_Lshift(cnp(n), obj.c())
	return obj2ObjErr(ret)
}

// Rshift returns the result of right shifting n by obj.  The equivalent Python
// is "n << obj".
//
// Return value: New Reference.
func (n *NumberProtocol) Rshift(obj *Base) (*Base, error) {
	ret := C.PyNumber_Rshift(cnp(n), obj.c())
	return obj2ObjErr(ret)
}

// And returns the bitwise and of n and obj.  The equivalent Python is "n &
// obj".
//
// Return value: New Reference.
func (n *NumberProtocol) And(obj *Base) (*Base, error) {
	ret := C.PyNumber_And(cnp(n), obj.c())
	return obj2ObjErr(ret)
}

// Xor returns the bitwise xor of n and obj.  The equivalent Python is "n ^
// obj".
//
// Return value: New Reference.
func (n *NumberProtocol) Xor(obj *Base) (*Base, error) {
	ret := C.PyNumber_Xor(cnp(n), obj.c())
	return obj2ObjErr(ret)
}

// Or returns the bitwise or of n and obj.  The equivalent Python is "n | obj".
//
// Return value: New Reference.
func (n *NumberProtocol) Or(obj *Base) (*Base, error) {
	ret := C.PyNumber_Or(cnp(n), obj.c())
	return obj2ObjErr(ret)
}

// InPlaceAdd returns the result of adding n and obj.  This is done in place if
// supported by n.  The equivalent Python is "n += obj".
//
// Return value: New Reference.
func (n *NumberProtocol) InPlaceAdd(obj *Base) (*Base, error) {
	ret := C.PyNumber_InPlaceAdd(cnp(n), obj.c())
	return obj2ObjErr(ret)
}

// InPlaceSubtract returns the result of subtracting obj from n.  This is done
// in place if supported by n.  The equivalent Python is "n -= obj".
//
// Return value: New Reference.
func (n *NumberProtocol) InPlaceSubtract(obj *Base) (*Base, error) {
	ret := C.PyNumber_InPlaceSubtract(cnp(n), obj.c())
	return obj2ObjErr(ret)
}

// InPlaceMultiply returns the result of multiplying n by obj.  This is done in
// place if supported by n.  The equivalent Python is "n *= obj".
//
// Return value: New Reference.
func (n *NumberProtocol) InPlaceMultiply(obj *Base) (*Base, error) {
	ret := C.PyNumber_InPlaceMultiply(cnp(n), obj.c())
	return obj2ObjErr(ret)
}

// InPlaceDivide returns the result of dividing n by obj.  This is done in place
// if supported by n.  The equivalent Python is "n /= obj".
//
// Return value: New Reference.
func (n *NumberProtocol) InPlaceDivide(obj *Base) (*Base, error) {
	ret := C.PyNumber_InPlaceDivide(cnp(n), obj.c())
	return obj2ObjErr(ret)
}

// TODO returns the ...
//
// Return value: New Reference.
func (n *NumberProtocol) InPlaceFloorDivide(obj *Base) (*Base, error) {
	ret := C.PyNumber_InPlaceFloorDivide(cnp(n), obj.c())
	return obj2ObjErr(ret)
}

// TODO returns the ...
//
// Return value: New Reference.
func (n *NumberProtocol) InPlaceTrueDivide(obj *Base) (*Base, error) {
	ret := C.PyNumber_InPlaceTrueDivide(cnp(n), obj.c())
	return obj2ObjErr(ret)
}

// InPlaceRemainder returns the remainder of n divided by obj.  This is done in
// place if supported by n.  The equivalent Python is "n %= obj".
//
// Return value: New Reference.
func (n *NumberProtocol) InPlaceRemainder(obj *Base) (*Base, error) {
	ret := C.PyNumber_InPlaceRemainder(cnp(n), obj.c())
	return obj2ObjErr(ret)
}

// InPlacePower returns the result of the Python "pow(n, obj1, obj2)".  This is
// done in place if supported by n.  If obj2 is None, then the Python "n **=
// obj" is also equivalent, if obj2 is not None, there is no equivalent in
// Python.
//
// Return value: New Reference.
func (n *NumberProtocol) InPlacePower(obj1, obj2 *Base) (*Base, error) {
	ret := C.PyNumber_InPlacePower(cnp(n), obj1.c(), obj2.c())
	return obj2ObjErr(ret)
}

// InPlaceLshift returns the result of left shifting n by obj.  This is done in
// place if supported by n.  The equivalent Python is "n <<= obj".
//
// Return value: New Reference.
func (n *NumberProtocol) InPlaceLshift(obj *Base) (*Base, error) {
	ret := C.PyNumber_InPlaceLshift(cnp(n), obj.c())
	return obj2ObjErr(ret)
}

// InPlaceRshift returns the result of right shifting n by obj.  This is done in
// place if supported by n.  The equivalent Python is "n >>= obj".
//
// Return value: New Reference.
func (n *NumberProtocol) InPlaceRshift(obj *Base) (*Base, error) {
	ret := C.PyNumber_InPlaceRshift(cnp(n), obj.c())
	return obj2ObjErr(ret)
}

// InPlaceAnd returns the bitwise and of n and obj.  This is done in place if
// supported by n.  The equivalent Python is "n &= obj".
//
// Return value: New Reference.
func (n *NumberProtocol) InPlaceAnd(obj *Base) (*Base, error) {
	ret := C.PyNumber_InPlaceAnd(cnp(n), obj.c())
	return obj2ObjErr(ret)
}

// InPlaceXor returns the bitwise xor of n and obj.  This is done in place if
// supported by n.  The equivalent Python is "n ^= obj".
//
// Return value: New Reference.
func (n *NumberProtocol) InPlaceXor(obj *Base) (*Base, error) {
	ret := C.PyNumber_InPlaceXor(cnp(n), obj.c())
	return obj2ObjErr(ret)
}

// InPlaceOr returns the bitwise or of n and obj.  This is done in place if
// supported by n.  The equivalent Python is "n |= obj".
//
// Return value: New Reference.
func (n *NumberProtocol) InPlaceOr(obj *Base) (*Base, error) {
	ret := C.PyNumber_InPlaceOr(cnp(n), obj.c())
	return obj2ObjErr(ret)
}

// PyNumber_Coerce: TODO

// PyNumber_CoerceEx: TODO

// PyNumber_Int: TODO

// PyNumber_Long: TODO

// PyNumber_Float: TODO

// PyNumber_Index: TODO

// PyNumber_ToBase: TODO

// PyNumber_AsSsize_t: TODO
