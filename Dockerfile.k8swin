FROM buildpack-deps:xenial

RUN apt-get update \
    && apt-get -y upgrade \
    && apt-get -y install apt-transport-https ca-certificates make gcc gcc-aarch64-linux-gnu rsync python-pip build-essential curl openssl vim jq \
    && rm -rf /var/lib/apt/lists/*

ENV GO_VERSION 1.8.3

RUN wget -q https://storage.googleapis.com/golang/go${GO_VERSION}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz && rm go${GO_VERSION}.linux-amd64.tar.gz

RUN curl -fsSL https://get.docker.com/ | sh

ENV GOPATH /gopath
ENV PATH "${PATH}:${GOPATH}/bin:/usr/local/go/bin"

RUN go get -u github.com/go-bindata/go-bindata/go-bindata

WORKDIR /gopath/src/k8s.io/kubernetes

ADD . /gopath/src/k8s.io/kubernetes
