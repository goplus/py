package py

// #include <Python.h>
// static inline void decref(PyObject *obj) { Py_DECREF(obj); }
import "C"
import "unsafe"

// ------------------------------------------------------------------------------------------

type GoModule struct {
	*Module
	Ctx RegisterCtx
}

func NewGoModule(name string, doc string, self interface{}) (mod GoModule, err error) {

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	var mdoc *C.char
	if doc != "" {
		mdoc = C.CString(doc)
		defer C.free(unsafe.Pointer(mdoc))
	}

	m := C.Py_InitModule4(cName, nil, mdoc, nil, C.PYTHON_API_VERSION)
	if m == nil {
		err = exception()
		return
	}

	mod.Module = (*Module)(unsafe.Pointer(m))
	mod.Ctx = Register(mod.Module.Dict(), name + ".", self)
	return
}

// ------------------------------------------------------------------------------------------

