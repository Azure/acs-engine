#!/bin/bash

docker run --rm -it -v "$PWD":/usr/src/acsengine -w /usr/src/acsengine golang bash

