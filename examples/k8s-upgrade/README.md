# Microsoft Azure Container Service Engine - Kubernetes Upgrade

## Overview

This document describes how to upgrade kubernetes version on a running cluster.
The upgrade is an **experimental** feature, and currently under development.

Supported scenario: upgrade from v1.5.3 to v1.6.2

The cluster definition file examples demonstrate initial cluster configurations:
- **kubernetes1.5.3.json** - Kubernetes cluster v1.5.3 with Linux agent pool
- **kubernetes1.5.3-win.json** - Kubernetes cluster v1.5.3 with Windows agent pool
- **kubernetes1.5.3-hybrid.json** - Kubernetes cluster v1.5.3 with Linux and Windows agent pools

The ***.env** files are used to set desired kubernetes version and instruct test framework to invoke post-deploy instructions implemented in **k8s-upgrade.sh** script.
