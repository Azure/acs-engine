FROM buildpack-deps:xenial

RUN apt-get update \
    && apt-get -y upgrade \
    && apt-get -y install python-pip make build-essential curl openssl vim jq \
    && rm -rf /var/lib/apt/lists/*

ENV GO_VERSION 1.8
RUN mkdir /tmp/godeb \
    && curl "https://godeb.s3.amazonaws.com/godeb-amd64.tar.gz" > /tmp/godeb/godeb.tar.gz \
    && (cd /tmp/godeb; tar zvxf godeb.tar.gz; ./godeb install "${GO_VERSION}") \
    && rm -rf /tmp/godeb

# See: https://github.com/Azure/azure-cli/blob/master/packaged_releases/bundled/README.md#using-the-bundled-installer
ENV AZURE_CLI_BUNDLE_VERSION 0.2.10-1
RUN mkdir /tmp/azurecli \
    && curl "https://azurecliprod.blob.core.windows.net/bundled/azure-cli_bundle_${AZURE_CLI_BUNDLE_VERSION}.tar.gz" > /tmp/azurecli/azure-cli_bundle.tar.gz \
    && (cd /tmp/azurecli \
      && tar -xvzf azure-cli_bundle.tar.gz \
      && azure-cli_bundle_*/installer --bin-dir /usr/local/bin) \
    && rm -rf /tmp/azurecli

RUN curl -fsSL https://get.docker.com/ | sh

ENV KUBECTL_VERSION 1.6.0
RUN curl "https://storage.googleapis.com/kubernetes-release/release/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl" > /usr/local/bin/kubectl \
    && chmod +x /usr/local/bin/kubectl

ENV GOPATH /gopath
ENV PATH "${PATH}:${GOPATH}/bin"

RUN git clone https://github.com/akesterson/cmdarg.git /tmp/cmdarg \
    && cd /tmp/cmdarg && make install && rm -rf /tmp/cmdarg
RUN git clone https://github.com/akesterson/shunit.git /tmp/shunit \
    && cd /tmp/shunit && make install && rm -rf /tmp/shunit

# Used by some CI jobs
ADD ./test/bootstrap/checkout-pr.sh /tmp/checkout-pr.sh

WORKDIR /gopath/src/github.com/Azure/acs-engine
