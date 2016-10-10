#!/bin/bash

docker run --rm -it -v "$PWD":/usr/src/acstgen -w /usr/src/acstgen golang bash

