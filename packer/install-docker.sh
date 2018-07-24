#!/bin/bash -eux

export DEBIAN_FRONTEND=noninteractive
apt_flags=(-o "Dpkg::Options::=--force-confnew" -qy)

DOCKER_REPO="https://apt.dockerproject.org/repo"
DOCKER_ENGINE_VERSION="1.13.*"

curl -fsSL https://aptdocker.azureedge.net/gpg > /tmp/aptdocker.gpg
apt-key add /tmp/aptdocker.gpg
echo "deb ${DOCKER_REPO} ubuntu-xenial main" | sudo tee /etc/apt/sources.list.d/docker.list
printf "Package: docker-engine\nPin: version ${DOCKER_ENGINE_VERSION}\nPin-Priority: 550\n" > /etc/apt/preferences.d/docker.pref

apt-get update -q
apt-get install "${apt_flags[@]}" docker-engine
