#!/usr/bin/env bash

set -eu -o pipefail
set -x

docker build --pull -t acs-engine .

docker run -it \
	--privileged \
	-v /var/run/docker.sock:/var/run/docker.sock \
	-v `pwd`:/gopath/src/github.com/Azure/acs-engine \
	-v ~/.azure:/root/.azure \
	-w /gopath/src/github.com/Azure/acs-engine \
		acs-engine /bin/bash

chown -R $(logname):$(id -gn $(logname)) . ~/.azure
