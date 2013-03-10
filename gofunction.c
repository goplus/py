#include <Python.h>
#include "_cgo_export.h"

int setMethod(PyMethodDef* d, int nin) {
	switch (nin) {
	case 3:
		d->ml_meth = (PyCFunction)goClassCallMethodKwds;
		d->ml_flags = METH_VARARGS | METH_KEYWORDS;
		break;
	case 2:
		d->ml_meth = (PyCFunction)goClassCallMethodArgs;
		d->ml_flags = METH_VARARGS;
		break;
	default:
		return -1;
	}
	return 0;
}

