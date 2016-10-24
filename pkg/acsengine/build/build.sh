#!/bin/bash

if ! which go-bindata ; then
	go get -u github.com/jteeuwen/go-bindata/...
fi

go generate
go build
