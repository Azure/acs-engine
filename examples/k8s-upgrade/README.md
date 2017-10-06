# Microsoft Azure Container Service Engine - Kubernetes Upgrade

## Overview

This document describes how to upgrade kubernetes version on a running cluster.

Supported scenarios:
- upgrade from v1.5 to v1.6
- upgrade from v1.6 to v1.7

The cluster definition file examples demonstrate initial cluster configurations:
- **kubernetes1.5.json** - Kubernetes cluster v1.5 with Linux agent pool
- **kubernetes1.5-win.json** - Kubernetes cluster v1.5 with Windows agent pool
- **kubernetes1.5-hybrid.json** - Kubernetes cluster v1.5 with Linux and Windows agent pools
- **kubernetes1.6.json** - Kubernetes cluster v1.6 with Linux agent pool

The ***.env** files are used to set desired kubernetes version and instruct test framework to invoke post-deploy instructions implemented in **k8s-upgrade.sh** script.
