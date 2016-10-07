#!/bin/bash

if [ ! -d generated ]
then
    mkdir -p generated
fi

TEMPLATE_NAME=$1



docker run --rm -it -v "$PWD":/usr/src/acstgen -w /usr/src/acstgen golang bash -c "go build && ./acstgen clusterdefinitions/$TEMPLATE_NAME > generated/$TEMPLATE_NAME"
