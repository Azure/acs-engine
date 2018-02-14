# Microsoft Azure Container Service Engine - Using GPUs with Kubernetes

If you created a Kubernetes cluster with one or multiple agent pool(s) whose VM size is `Standard_NC*` or `Standard_NV*` you can schedule GPU workload on your cluster.
The NVIDIA drivers are automatically installed on every GPU agent in your cluster, so you don't need to do that manually, unless you require a specific version of the drivers. Currently, the installed driver is version 378.13.

To make sure everything is fine, run `kubectl describe node <name-of-a-gpu-node>`. You should see the correct number of GPU reported (in this example shows 2 GPU for a NC12 VM):
```
[...]
Capacity:
 alpha.kubernetes.io/nvidia-gpu:	2
 cpu:					12
[...]
```

If `alpha.kubernetes.io/nvidia-gpu` is `0` and you just created the cluster, you might have to wait a little bit. The driver installation takes about 12 minutes, and the node might join the cluster before the installation is completed. After a few minute the node should restart, and report the correct number of GPUs.

## Running a GPU-enabled container

When running a GPU container, you will need to specify how many GPU you want to use. If you don't specify a GPU count, kubernetes will asumme you don't require any, and will not map the device into the container.
You will also need to mount the drivers from the host (the kubernetes agent) into the container.

On the host, the drivers are installed under `/usr/lib/nvidia-384`.

Here is an example template running TensorFlow: 

```
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  labels:
    app: tensorflow
  name: tensorflow
spec:
  template:
    metadata:
      labels:
        app: tensorflow
    spec:     
      containers:
      - name: tensorflow
        image: <SOME_IMAGE>
        command: <SOME_COMMAND>    
        imagePullPolicy: IfNotPresent
        resources:
          limits:
            alpha.kubernetes.io/nvidia-gpu: 1
        volumeMounts:
        - mountPath: /usr/local/nvidia/bin
          name: bin
        - mountPath: /usr/local/nvidia/lib64
          name: lib
        - mountPath: /usr/lib/x86_64-linux-gnu/libcuda.so.1
          name: libcuda
      volumes:
        - name: bin
          hostPath: 
            path: /usr/lib/nvidia-384/bin
        - name: lib
          hostPath: 
            path: /usr/lib/nvidia-384
        - name: libcuda
          hostPath:
            path: /usr/lib/x86_64-linux-gnu/libcuda.so.1
```

We specify `alpha.kubernetes.io/nvidia-gpu: 1` in the resources limits, and we mount the drivers from the host into the container.

Some libraries, such as `libcuda.so` are installed under `/usr/lib/x86_64-linux-gnu` on the host, you might need to mount them separatly as shown above based on your needs.



