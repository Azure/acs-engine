#!/bin/bash

if ! which go-bindata ; then
	go get -u github.com/lestrrat/go-bindata/...
fi

go generate
go build
