package py

/*
#include <Python.h>

static inline int codeCheck(PyObject *o) { return PyCode_Check(o); }
static inline void decref(PyObject *obj) { Py_DECREF(obj); }

static inline FILE* openFile(char* name) {
	return fopen(name, "r");
}

static inline PyObject* compileString(char* text, char* filename, int start) {
	return Py_CompileString(text, filename, start);
}

static PyObject* compileFile(FILE* f, char* name, int start) {
	struct _node *n = PyParser_SimpleParseFile(f, name, start);
	if (!n) return NULL;
	return (PyObject*)PyNode_Compile(n, name);
}
*/
import "C"
import "unsafe"

// ------------------------------------------------------------------------------------------
// type StartToken

type StartToken int

const (
	EvalInput = StartToken(C.Py_eval_input) // for isolated expressions

	FileInput = StartToken(C.Py_file_input) // for sequences of statements as read from a file or other source;
	// to use when compiling arbitrarily long Python source code.

	SingleInput = StartToken(C.Py_single_input) // for a single statement; used for the interactive interpreter loop.
)

// ------------------------------------------------------------------------------------------
// type Code

type Code struct {
	Base
	o C.PyCodeObject
}

// CodeType is the Type object that represents the Code type.
var CodeType = (*Type)(unsafe.Pointer(&C.PyCode_Type))

func newCode(obj *C.PyObject) *Code {
	return (*Code)(unsafe.Pointer(obj))
}

func AsCode(o *Base) (v *Code, ok bool) {
	if ok = C.codeCheck(o.c()) != 0; ok {
		v = newCode(o.c())
	}
	return
}

func Compile(text, filename string, start StartToken) (*Code, error) {
	t := C.CString(text)
	defer C.free(unsafe.Pointer(t))

	fn := C.CString(filename)
	defer C.free(unsafe.Pointer(fn))

	ret := C.compileString(t, fn, C.int(start))
	if ret == nil {
		return nil, exception()
	}
	return newCode(ret), nil
}

func CompileFile(name string, start StartToken) (*Code, error) {
	fn := C.CString(name)
	defer C.free(unsafe.Pointer(fn))

	file, err := C.openFile(fn)
	if file == nil {
		return nil, err
	}
	defer C.fclose(file)

	ret := C.compileFile(file, fn, C.int(start))
	if ret == nil {
		return nil, exception()
	}
	return newCode(ret), nil
}

// Return value: New reference.
func (code *Code) Eval(globals, locals *Base) (*Base, error) {
	pyCode := (*C.PyCodeObject)(unsafe.Pointer(code))
	ret := C.PyEval_EvalCode(pyCode, globals.c(), locals.c())
	return obj2ObjErr(ret)
}

func (code *Code) Run(globals, locals *Base) error {
	pyCode := (*C.PyCodeObject)(unsafe.Pointer(code))
	ret := C.PyEval_EvalCode(pyCode, globals.c(), locals.c())
	if ret == nil {
		return exception()
	}
	C.decref(ret)
	return nil
}

// ------------------------------------------------------------------------------------------

func Run(text string) error {

	t := C.CString(text)
	defer C.free(unsafe.Pointer(t))

	ret := C.PyRun_SimpleStringFlags(t, nil)
	return int2Err(ret)
}

// Return a dictionary of the builtins in the current execution frame, or the interpreter of
// the thread state if no frame is currently executing.
//
// Return value: Borrowed reference.
func GetBuiltins() *Base {
	ret := C.PyEval_GetBuiltins()
	return newObject(ret)
}

// Return a dictionary of the global variables in the current execution frame,
// or NULL if no frame is currently executing.
//
// Return value: Borrowed reference
func GetLocals() *Base {
	ret := C.PyEval_GetLocals()
	return newObject(ret)
}

// Return a dictionary of the local variables in the current execution frame,
// or NULL if no frame is currently executing.
//
// Return value: Borrowed reference
func GetGlobals() *Base {
	ret := C.PyEval_GetGlobals()
	return newObject(ret)
}

// ------------------------------------------------------------------------------------------
