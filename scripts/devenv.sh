#!/usr/bin/env bash

set -eu -o pipefail
set -x

docker build --pull -t \
	--build-arg http_proxy=$http_proxy \
	--build-arg https_proxy=$https_proxy \
	--build-arg no_proxy=$no_proxy \
	--build-arg HTTP_PROXY=$HTTP_PROXY \
	--build-arg HTTPS_PROXY=$HTTPS_PROXY \
	--build-arg NO_PROXY=$NO_PROXY \
	acs-engine .

docker run -it \
	--privileged \
	-v /var/run/docker.sock:/var/run/docker.sock \
	-v `pwd`:/gopath/src/github.com/Azure/acs-engine \
	-v ~/.azure:/root/.azure \
	-w /gopath/src/github.com/Azure/acs-engine \
		acs-engine /bin/bash

chown -R "$(logname):$(id -gn $(logname))" . ~/.azure
