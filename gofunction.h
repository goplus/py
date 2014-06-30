#ifndef QBOX_GOPY_GOFUNCTION_H
#define QBOX_GOPY_GOFUNCTION_H

int setMethod(PyMethodDef* d, int nin);
static inline void decref(PyObject *obj) { Py_DECREF(obj); }

#endif /* _GO_PYTHON_UTILS_H */

