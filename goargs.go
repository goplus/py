package py

import (
	"github.com/qiniu/x/errors"
	"github.com/qiniu/x/log"
	"reflect"
	"syscall"
)

// ------------------------------------------------------------------------------------------

func ToInt(in *Base) (int, bool) {

	l, ok := ToLong(in)
	return int(l), ok
}

func ToLong(in *Base) (int64, bool) {

	if v, ok := AsInt(in); ok {
		return int64(v.Int()), true
	}
	if v, ok := AsLong(in); ok {
		return int64(v.Long()), true
	}
	return 0, false
}

func ToString(in *Base) (string, bool) {

	if v, ok := AsString(in); ok {
		return v.String(), true
	}
	return "", false
}

func ToInterface(in *Base) (v interface{}, ok bool) {

	if v, ok = ToLong(in); ok {
		return
	} else if v, ok = ToString(in); ok {
		return
	}
	return
}

// ------------------------------------------------------------------------------------------

func assignToMap(in *Base, out reflect.Value) (err error) {

	dict, ok := AsDict(in)
	if !ok {
		err = errors.Info(syscall.EINVAL, "py.AssignTo", "uncompatible type")
		return
	}

	mapTy := out.Type()
	m := reflect.MakeMap(mapTy)

	keyTy := mapTy.Key()
	valTy := mapTy.Elem()

	log.Debug("assignToMap:", mapTy, keyTy, valTy)

	var iter DictIter
	var k, v *Base
	for dict.Next(&iter, &k, &v) {
		kout := reflect.New(keyTy)
		err = AssignTo(k, kout.Interface())
		if err != nil {
			err = errors.Info(err, "py.AssignTo", "assign map key").Detail(err)
			return
		}
		vout := reflect.New(valTy)
		err = AssignTo(v, vout.Interface())
		if err != nil {
			err = errors.Info(err, "py.AssignTo", "assign map val").Detail(err)
			return
		}
		m.SetMapIndex(kout.Elem(), vout.Elem())
	}

	out.Set(m)
	return nil
}

func assignToComplex(in *Base, out1 reflect.Value) (err error) {

	if out1.Kind() != reflect.Ptr {
		err = errors.Info(syscall.EINVAL, "py.AssignTo", "not assignable")
		return
	}

	out := out1.Elem()
	switch out.Kind() {
	case reflect.Map:
		return assignToMap(in, out)
	default:
		err = errors.Info(syscall.EINVAL, "py.AssignTo", "unsupported input type", reflect.TypeOf(out))
		return
	}
	return
}

func AssignTo(in *Base, out interface{}) (err error) {

	var ok bool
	switch v := out.(type) {
	case *string:
		*v, ok = ToString(in)
	case *int64:
		*v, ok = ToLong(in)
	case *int:
		*v, ok = ToInt(in)
	case *interface{}:
		*v, ok = ToInterface(in)
		return
	default:
		return assignToComplex(in, reflect.ValueOf(out))
	}
	if !ok {
		err = errors.Info(syscall.EINVAL, "py.AssignTo", "can not convert type", reflect.TypeOf(out))
	}
	return
}

// ------------------------------------------------------------------------------------------

func Parse(in *Tuple, out ...interface{}) (err error) {

	n := in.Size()
	if n != len(out) {
		err = errors.Info(syscall.EINVAL, "py.Parse", "invalid argument count")
		return
	}

	for i := 0; i < n; i++ {
		v2, err2 := in.GetItem(i)
		if err2 != nil {
			err = errors.Info(err2, "py.Parse", "invalid argument", i+1).Detail(err2)
			return
		}
		err2 = AssignTo(v2, out[i])
		if err2 != nil {
			err = errors.Info(err2, "py.Parse", "assign argument failed", i+1).Detail(err2)
			return
		}
	}
	return
}

func ParseV(in *Tuple, out ...interface{}) (err error) {

	n1 := in.Size()
	n := len(out) - 1
	if n1 < n || n < 0 {
		err = errors.Info(syscall.EINVAL, "py.ParseV", "argument count is not enough")
		return
	}

	slicePtr := reflect.TypeOf(out[n])
	if slicePtr.Kind() != reflect.Ptr {
		err = errors.Info(syscall.EINVAL, "py.ParseV", "last argument is not a slice pointer")
		return
	}

	sliceTy := slicePtr.Elem()
	if sliceTy.Kind() != reflect.Slice {
		err = errors.Info(syscall.EINVAL, "py.ParseV", "last argument is not a slice pointer")
		return
	}

	for i := 0; i < n; i++ {
		v2, err2 := in.GetItem(i)
		if err2 != nil {
			err = errors.Info(err2, "py.ParseV", "invalid argument", i+1).Detail(err2)
			return
		}
		err2 = AssignTo(v2, out[i])
		if err2 != nil {
			err = errors.Info(err2, "py.ParseV", "assign argument failed", i+1).Detail(err2)
			return
		}
	}

	slice := reflect.MakeSlice(sliceTy, n1-n, n1-n)
	for i := n; i < n1; i++ {
		v2, err2 := in.GetItem(i)
		if err2 != nil {
			err = errors.Info(err2, "py.ParseV", "invalid argument", i+1).Detail(err2)
			return
		}
		err2 = AssignTo(v2, slice.Index(i-n).Addr().Interface())
		if err2 != nil {
			err = errors.Info(err2, "py.ParseV", "assign argument failed", i+1).Detail(err2)
			return
		}
	}
	reflect.ValueOf(out[n]).Elem().Set(slice)
	return
}

// ------------------------------------------------------------------------------------------
