Golang bindings to the CPython C-API
=====

[![Build Status](https://travis-ci.org/qiniu/py.png?branch=master)](https://travis-ci.org/qiniu/py) [![GoDoc](https://godoc.org/github.com/qiniu/py?status.svg)](https://godoc.org/github.com/qiniu/py)

[![Qiniu Logo](http://open.qiniudn.com/logo.png)](http://www.qiniu.com/)

py is Golang bindings to the CPython C-API.

py project's homepage is: https://github.com/qiniu/py

# Document

See http://godoc.org/github.com/qiniu/py

# Install

```
go get github.com/qiniu/py
```

# Example

```{go}
package main

import (
	"fmt"
	"github.com/qiniu/x/log"
	"github.com/qiniu/py"
)

// -------------------------------------------------------------------

type FooModule struct {
}

func (r *FooModule) Py_bar(args *py.Tuple) (ret *py.Base, err error) {
	var i int
	var s string
	err = py.Parse(args, &i, &s)
	if err != nil {
		return
	}
	fmt.Println("call foo.bar:", i, s)
	return py.IncNone(), nil
}

func (r *FooModule) Py_bar2(args *py.Tuple) (ret *py.Base, err error) {
	var i int
	var s []string
	err = py.ParseV(args, &i, &s)
	if err != nil {
		return
	}
	fmt.Println("call foo.bar2:", i, s)
	return py.IncNone(), nil
}

// -------------------------------------------------------------------

const pyCode = `

import foo
foo.bar(1, 'Hello')
foo.bar2(1, 'Hello', 'world!')
`

func main() {

	gomod, err := py.NewGoModule("foo", "", new(FooModule))
	if err != nil {
		log.Fatal("NewGoModule failed:", err)
	}
	defer gomod.Decref()

	code, err := py.Compile(pyCode, "", py.FileInput)
	if err != nil {
		log.Fatal("Compile failed:", err)
	}
	defer code.Decref()

	mod, err := py.ExecCodeModule("test", code.Obj())
	if err != nil {
		log.Fatal("ExecCodeModule failed:", err)
	}
	defer mod.Decref()
}

// -------------------------------------------------------------------
```

