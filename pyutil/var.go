package pyutil

import (
	"strings"
	"reflect"
	"github.com/qiniu/py"
)

// ------------------------------------------------------------------------------------------

type Config struct {
	Cate string
	SliceAsList bool
}

var DefaultConfig = &Config {
	Cate: "json",
}

// ------------------------------------------------------------------------------------------

func tagName(tag string) (string) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx]
	}
	return tag
}

func newStruct(sv reflect.Value, cfg *Config) (ret *py.Base, ok bool) {

	dict := py.NewDict()

	st := sv.Type()
	for i := 0; i < sv.NumField(); i++ {
		sf := st.Field(i)
		tag := sf.Tag.Get(cfg.Cate)
		if tag == "" {
			return nil, false
		}
		name := tagName(tag)
		val := sv.Field(i)
		val1, ok1 := NewVarEx(val.Interface(), cfg)
		if !ok1 {
			dict.Decref()
			return nil, false
		}
		dict.SetItemString(name, val1)
		val1.Decref()
	}

	return dict.Obj(), true
}

func newMap(v reflect.Value, cfg *Config) (ret *py.Base, ok bool) {

	dict := py.NewDict()

	keys := v.MapKeys()
	for _, key := range keys {
		key1, ok1 := NewVarEx(key.Interface(), cfg)
		if !ok1 {
			dict.Decref()
			return nil, false
		}
		val1, ok1 := NewVarEx(v.MapIndex(key).Interface(), cfg)
		if !ok1 {
			key1.Decref()
			dict.Decref()
			return nil, false
		}
		dict.SetItem(key1, val1)
		key1.Decref()
		val1.Decref()
	}

	return dict.Obj(), true
}

func newComplex(val reflect.Value, cfg *Config) (ret *py.Base, ok bool) {

retry:
	switch val.Kind() {
	case reflect.Struct:
		return newStruct(val, cfg)
	case reflect.Map:
		return newMap(val, cfg)
	case reflect.Ptr, reflect.Interface:
		val = val.Elem()
		goto retry
	}
	return nil, false
}

// ------------------------------------------------------------------------------------------

func NewVarEx(val interface{}, cfg *Config) (ret *py.Base, ok bool) {

	switch v := val.(type) {
	case int:
		return py.NewInt(v).Obj(), true
	case int64:
		return py.NewLong(v).Obj(), true
	case string:
		return py.NewString(v).Obj(), true
	}
	return newComplex(reflect.ValueOf(val), cfg)
}

func NewVar(val interface{}) (ret *py.Base, ok bool) {

	return NewVarEx(val, DefaultConfig)
}

// ------------------------------------------------------------------------------------------
