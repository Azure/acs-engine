FROM alpine:3.7

ARG ACSENGINE_VERSION=0.16.0
ARG BUILD_DATE

# Metadata as defined at http://label-schema.org
LABEL maintainer="Microsoft" \
      org.label-schema.schema-version="1.0" \
      org.label-schema.vendor="Microsoft" \
      org.label-schema.name="Azure Container Service Engine (acs-engine)" \
      org.label-schema.version=$ACSENGINE_VERSION \
      org.label-schema.license="MIT" \
      org.label-schema.description="The Azure Container Service Engine (acs-engine) generates ARM (Azure Resource Manager) templates for Docker enabled clusters on Microsoft Azure with your choice of DCOS, Kubernetes, or Swarm orchestrators." \
      org.label-schema.url="https://github.com/Azure/acs-engine" \
      org.label-schema.usage="https://github.com/Azure/acs-engine/blob/master/docs/acsengine.md" \
      org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.vcs-url="https://github.com/Azure/acs-engine.git" \
      org.label-schema.docker.cmd="docker run -v \${PWD}:/acs-engine/workspace -it --rm microsoft/acs-engine:$ACSENGINE_VERSION"

RUN apk add --no-cache ca-certificates

ADD "https://github.com/Azure/acs-engine/releases/download/v${ACSENGINE_VERSION}/acs-engine-v${ACSENGINE_VERSION}-linux-amd64.tar.gz" /tmp/acs-engine.tgz

RUN mkdir /opt/ && \
    tar xvzf /tmp/acs-engine.tgz -C /opt/ && \
    rm /tmp/acs-engine.tgz && \
    chown -R root:root /opt/acs-engine-v${ACSENGINE_VERSION}-linux-amd64 && \
    ln -s /opt/acs-engine-v${ACSENGINE_VERSION}-linux-amd64/acs-engine /usr/local/bin/acs-engine && \
    chmod +x /usr/local/bin/acs-engine

WORKDIR /acs-engine/workspace

ENTRYPOINT [ "acs-engine" ]
CMD [ "--help" ]
