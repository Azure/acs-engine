# Microsoft Azure Container Service Engine - Kubernetes Multi-GPU support Walkthrough

## Deployment

Here are the steps to deploy a simple Kubernetes cluster with multi-GPU support:

1. [Install a Kubernetes cluster][Kubernetes Walkthrough](kubernetes.md) - shows how to create a Kubernetes cluster.
  > NOTE: Make sure to configure the agent nodes with vm size `Standard_NC12` or above to utilize the GPUs

2. Install drivers:
  * SSH into each node and run the following scripts : 
  install-nvidia-driver.sh
  ```
  curl -L -sf https://raw.githubusercontent.com/ritazh/acs-k8s-gpu/master/install-nvidia-driver.sh | sudo sh
  ```

  To verify, when you run `kubectl describe node <node-name>`, you should get something like the following:

  ```
  Capacity:
  alpha.kubernetes.io/nvidia-gpu:    2
  cpu:                               12
  memory:                            115505744Ki
  pods:                              110
  ```

3. Scheduling a multi-GPU container

* You need to specify `alpha.kubernetes.io/nvidia-gpu: 2` as a limit
* You need to expose the drivers to the container as a volume. If you are using TF original docker image, it is based on ubuntu 16.04, just like your cluster's VM, so you can just mount `/usr/bin` and `/usr/lib/x86_64-linux-gnu`, it's a bit dirty but it works. Ideally, improve the previous script to install the driver in a specific directory and only expose this one.

``` yaml
apiVersion: v1
kind: Pod
metadata:
  name: gpu-test
spec:
  volumes:
  - name: binaries
    hostPath:
      path: /usr/bin/
  - name: libraries
    hostPath:
      path: /usr/lib/x86_64-linux-gnu
  containers:
  - name: tensorflow
    image: gcr.io/tensorflow/tensorflow:latest-gpu
    ports:
    - containerPort: 8888
    resources:
      limits:
        alpha.kubernetes.io/nvidia-gpu: 2
    volumeMounts:
    - mountPath: /usr/bin/
      name: binaries
    - mountPath: /usr/lib/x86_64-linux-gnu
      name: libraries
```
To verify, when you run `kubectl describe pod <pod-name>`, you see get the following:

```
Successfully assigned gpu-test to k8s-agentpool1-10960440-1
```
