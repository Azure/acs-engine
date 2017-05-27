FROM quay.io/deis/go-dev:v0.26.0

RUN apt-get update \
    && apt-get install -y openssh-client python-pip \
    && pip install --disable-pip-version-check --no-cache-dir azure-cli==2.0.10