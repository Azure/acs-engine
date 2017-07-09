#!/bin/bash

/bin/systemctl stop kubelet-restart.service
/bin/systemctl stop kubelet-restart.timer || true
