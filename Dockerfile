FROM buildpack-deps:xenial

RUN apt-get update \
	&& apt-get -y upgrade \
	&& apt-get -y install python-pip make build-essential curl openssl vim jq \
	&& rm -rf /var/lib/apt/lists/*

RUN curl 'https://godeb.s3.amazonaws.com/godeb-amd64.tar.gz' > /tmp/godeb.tar.gz \
	&& (cd /tmp/; tar zvxf godeb.tar.gz; ./godeb install 1.7.3)

RUN pip install azure-cli==0.1.0b8

RUN curl https://storage.googleapis.com/kubernetes-release/release/v1.4.4/bin/linux/amd64/kubectl > /usr/local/bin/kubectl \
	&& chmod +x /usr/local/bin/kubectl

ENV GOPATH /gopath
ENV PATH "${PATH}:${GOPATH}/bin"
RUN go get -u github.com/golang/lint/golint
RUN go get -u github.com/jteeuwen/go-bindata/...
