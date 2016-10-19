#!/bin/bash

if [ ! -d generated ]
then
    mkdir -p generated
fi

TEMPLATE_NAME=$1



docker run --rm -it -v "$PWD":/usr/src/acsengine -w /usr/src/acsengine golang bash -c "go build && ./acsengine clusterdefinitions/$TEMPLATE_NAME > generated/$TEMPLATE_NAME"
