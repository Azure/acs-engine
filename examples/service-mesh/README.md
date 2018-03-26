# Kubernetes - Service Mesh

There are numerous implementations of a service mesh which integrate with kubernetes such as Istio, [Linkerd](http://linkerd.io), and [Conduit](https://conduit.io/).  [This is one blog post](https://medium.com/microservices-in-practice/service-mesh-for-microservices-2953109a3c9a) which explains some fundamentals behind what it is and why to use it.

Some service mesh implementations **may** benefit from or require additional [customizations to the kubernetes cluster itself](https://github.com/Azure/acs-engine/blob/master/docs/clusterdefinition.md).

## Istio

The `istio.json` file in this directory enables the kubernetes API server options to support automatic sidecar injection using [Isitio](https://istio.io/).  If automatic sidecar injection isn't enabled, then all services must then manually inject the sidecar configuration into every deployment, every time.  

The main changes this configuration makes is adding these flags to the apiserver `Initializers,MutatingAdmissionWebhook,ValidatingAdmissionWebhook` and starting using the `runtime-config` with `admissionregistration.k8s.io/v1alpha1`.

> Note: The default acs-engine apiserver options `AlwaysPullImages` and `SecurityContextDeny` were removed from this configuration in order to have the Istio book info examples work without any errors.  Consider enabling these for a production cluster.


### Post installation

Once the template has been successfully deployed, then Istio can be installed via either:

1. Manual - follow the website [Installation steps](https://istio.io/docs/setup/kubernetes/quick-start.html#installation-steps).
1. Helm Chart - is maintained in the Istio repository itself (no longer hub.kubeapps.com).  [See these instructions on the Istio website](https://istio.io/docs/setup/kubernetes/helm.html).

> Note: So far it seems the manual steps are more well maintained and up-to-date than the helm chart.

After Istio has been installed, consider [walking through the various Tasks](https://istio.io/docs/tasks/) which use the Book info example application.