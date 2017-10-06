package armhelpers

import (
	"time"

	log "github.com/sirupsen/logrus"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	policy "k8s.io/client-go/pkg/apis/policy/v1beta1"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	evictionKind        = "Eviction"
	evictionSubresource = "pods/eviction"
)

//KubernetesClientSetClient kubernetes client hooked up to a live api server
type KubernetesClientSetClient struct {
	clientset         *kubernetes.Clientset
	interval, timeout time.Duration
}

//GetKubernetesClient returns a KubernetesClient hooked up to the api server at the masterURL
func (az *AzureClient) GetKubernetesClient(masterURL, kubeConfig string, interval, timeout time.Duration) (KubernetesClient, error) {
	// creates the clientset
	config, err := clientcmd.BuildConfigFromKubeconfigGetter(masterURL, func() (*clientcmdapi.Config, error) { return clientcmd.Load([]byte(kubeConfig)) })
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &KubernetesClientSetClient{clientset: clientset, interval: interval, timeout: timeout}, nil
}

//ListPods returns all Pods running on the passed in node
func (c *KubernetesClientSetClient) ListPods(node *v1.Node) (*v1.PodList, error) {
	return c.clientset.Core().Pods(metav1.NamespaceAll).List(metav1.ListOptions{
		FieldSelector: fields.SelectorFromSet(fields.Set{"spec.nodeName": node.Name}).String()})
}

//GetNode returns details about node with passed in name
func (c *KubernetesClientSetClient) GetNode(name string) (*v1.Node, error) {
	return c.clientset.Nodes().Get(name, metav1.GetOptions{})
}

//UpdateNode updates the node in the api server with the passed in info
func (c *KubernetesClientSetClient) UpdateNode(node *v1.Node) (*v1.Node, error) {
	return c.clientset.Nodes().Update(node)
}

//SupportEviction queries the api server to discover if it supports eviction, and returns supported type if it is supported
func (c *KubernetesClientSetClient) SupportEviction() (string, error) {
	discoveryClient := c.clientset.Discovery()
	groupList, err := discoveryClient.ServerGroups()
	if err != nil {
		return "", err
	}
	foundPolicyGroup := false
	var policyGroupVersion string
	for _, group := range groupList.Groups {
		if group.Name == "policy" {
			foundPolicyGroup = true
			policyGroupVersion = group.PreferredVersion.GroupVersion
			break
		}
	}
	if !foundPolicyGroup {
		return "", nil
	}
	resourceList, err := discoveryClient.ServerResourcesForGroupVersion("v1")
	if err != nil {
		return "", err
	}
	for _, resource := range resourceList.APIResources {
		if resource.Name == evictionSubresource && resource.Kind == evictionKind {
			return policyGroupVersion, nil
		}
	}
	return "", nil
}

//DeletePod deletes the passed in pod
func (c *KubernetesClientSetClient) DeletePod(pod *v1.Pod) error {
	return c.clientset.Core().Pods(pod.Namespace).Delete(pod.Name, &metav1.DeleteOptions{})
}

//EvictPod evicts the passed in pod using the passed in api version
func (c *KubernetesClientSetClient) EvictPod(pod *v1.Pod, policyGroupVersion string) error {
	eviction := &policy.Eviction{
		TypeMeta: metav1.TypeMeta{
			APIVersion: policyGroupVersion,
			Kind:       evictionKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      pod.Name,
			Namespace: pod.Namespace,
		},
	}
	return c.clientset.Policy().Evictions(eviction.Namespace).Evict(eviction)
}

//GetPod returns the pod
func (c *KubernetesClientSetClient) getPod(namespace, name string) (*v1.Pod, error) {
	return c.clientset.Core().Pods(namespace).Get(name, metav1.GetOptions{})
}

//WaitForDelete waits until all pods are deleted. Returns all pods not deleted and an error on failure
func (c *KubernetesClientSetClient) WaitForDelete(logger *log.Entry, pods []v1.Pod, usingEviction bool) ([]v1.Pod, error) {
	var verbStr string
	if usingEviction {
		verbStr = "evicted"
	} else {
		verbStr = "deleted"
	}
	err := wait.PollImmediate(c.interval, c.timeout, func() (bool, error) {
		pendingPods := []v1.Pod{}
		for i, pod := range pods {
			p, err := c.getPod(pod.Namespace, pod.Name)
			if apierrors.IsNotFound(err) || (p != nil && p.ObjectMeta.UID != pod.ObjectMeta.UID) {
				logger.Infof("%s pod successfully %s", pod.Name, verbStr)
				continue
			} else if err != nil {
				return false, err
			} else {
				pendingPods = append(pendingPods, pods[i])
			}
		}
		pods = pendingPods
		if len(pendingPods) > 0 {
			return false, nil
		}
		return true, nil
	})
	return pods, err
}
