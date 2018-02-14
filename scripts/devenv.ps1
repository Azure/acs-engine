$pwd = (Get-Location).Path

docker build --pull -t acs-engine .
docker run --security-opt seccomp:unconfined -it `
	-v ${pwd}:/gopath/src/github.com/Azure/acs-engine `
	-w /gopath/src/github.com/Azure/acs-engine `
		acs-engine /bin/bash

