#!/usr/bin/env bash

set -eu -o pipefail
set -x

docker build -t microsoft/acs-engine .

docker run -it \
	--privileged \
	-v /var/run/docker.sock:/var/run/docker.sock \
	-v ~/.azure:/root/.azure \
	-w /gopath/src/github.com/Azure/acs-engine \
		acs-engine /bin/bash

chown -R "$(logname):$(id -gn $(logname))" . ~/.azure
