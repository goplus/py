#!/bin/bash
export CGO_CFLAGS="-I/usr/include/python2.7"
export CGO_LDFLAGS="-lpython2.7"
go install
