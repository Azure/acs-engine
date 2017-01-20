FROM buildpack-deps:xenial

ENV GO_VERSION 1.7.3
ENV KUBECTL_VERSION 1.4.4
ENV AZURE_CLI_VERSION 0.1.0b8

RUN apt-get update \
    && apt-get -y upgrade \
    && apt-get -y install python-pip make build-essential curl openssl vim jq \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir /tmp/godeb \
    && curl "https://godeb.s3.amazonaws.com/godeb-amd64.tar.gz" > /tmp/godeb/godeb.tar.gz \
    && (cd /tmp/godeb; tar zvxf godeb.tar.gz; ./godeb install "${GO_VERSION}") \
    && rm -rf /tmp/godeb

RUN pip install "azure-cli==${AZURE_CLI_VERSION}"

RUN curl "https://storage.googleapis.com/kubernetes-release/release/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl" > /usr/local/bin/kubectl \
    && chmod +x /usr/local/bin/kubectl

ENV GOPATH /gopath
ENV PATH "${PATH}:${GOPATH}/bin"
RUN go get -u github.com/golang/lint/golint
RUN go get -u github.com/ghodss/yaml
RUN go get -u github.com/jteeuwen/go-bindata/...
