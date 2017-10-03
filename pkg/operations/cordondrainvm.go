package operations

import (
	"fmt"
	"math"
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
	EvictionKind        = "Eviction"
	EvictionSubresource = "pods/eviction"
	interval            = time.Second * 1
	MirrorPodAnnotation = "kubernetes.io/config.mirror"
)

type drainOperation struct {
	clientset *kubernetes.Clientset
	node      *v1.Node
	logger    *log.Entry
	timeout   time.Duration
}

type podFilter func(v1.Pod) bool

// SafelyDrainNode safely drains a node so that it can be deleted from the cluster
func SafelyDrainNode(logger *log.Entry, masterUrl, kubeConfig, nodeName string) error {
	//get client using kubeconfig
	clientset, err := newInClusterKubeClient(masterUrl, kubeConfig)
	if err != nil {
		return err
	}

	//Mark the node unschedulable
	node, err := clientset.Nodes().Get(nodeName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	node.Spec.Unschedulable = true
	node, err = clientset.Nodes().Update(node)
	if err != nil {
		return err
	}

	//Evict pods in node
	drainOp := &drainOperation{clientset: clientset, node: node, logger: logger}
	drainOp.deleteOrEvictPodsSimple()

	return nil
}

func newInClusterKubeClient(masterUrl, kubeConfig string) (*kubernetes.Clientset, error) {
	// creates the clientset
	config, err := clientcmd.BuildConfigFromKubeconfigGetter(masterUrl, func() (*clientcmdapi.Config, error) { return clientcmd.Load([]byte(kubeConfig)) })
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func (o *drainOperation) deleteOrEvictPodsSimple() error {
	pods, err := o.getPodsForDeletion()
	if err != nil {
		return err
	}

	err = o.deleteOrEvictPods(pods)
	if err != nil {
		pendingPods, newErr := o.getPodsForDeletion()
		if newErr != nil {
			return newErr
		}
		o.logger.Errorf("There are pending pods when an error occurred: %v\n", err)
		for _, pendingPod := range pendingPods {
			o.logger.Errorf("%s/%s\n", "pod", pendingPod.Name)
		}
	}
	return err
}

func mirrorPodFilter(pod v1.Pod) bool {
	if _, found := pod.ObjectMeta.Annotations[MirrorPodAnnotation]; found {
		return false
	}
	return true
}

// getPodsForDeletion returns all the pods we're going to delete.  If there are
// any pods preventing us from deleting, we return that list in an error.
func (o *drainOperation) getPodsForDeletion() (pods []v1.Pod, err error) {
	podList, err := o.clientset.Core().Pods(metav1.NamespaceAll).List(metav1.ListOptions{
		FieldSelector: fields.SelectorFromSet(fields.Set{"spec.nodeName": o.node.Name}).String()})
	if err != nil {
		return pods, err
	}

	for _, pod := range podList.Items {
		podOk := true
		for _, filt := range []podFilter{
			mirrorPodFilter,
			// localStorageFilter,
			//unreplicatedFilter,
		} {
			podOk = podOk && filt(pod)
		}
		if podOk {
			pods = append(pods, pod)
		}
	}
	return pods, nil
}

func (o *drainOperation) deletePod(pod v1.Pod) error {
	return o.clientset.Core().Pods(pod.Namespace).Delete(pod.Name, &metav1.DeleteOptions{})
}

func (o *drainOperation) evictPod(pod v1.Pod, policyGroupVersion string) error {
	eviction := &policy.Eviction{
		TypeMeta: metav1.TypeMeta{
			APIVersion: policyGroupVersion,
			Kind:       EvictionKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      pod.Name,
			Namespace: pod.Namespace,
		},
	}
	// Remember to change change the URL manipulation func when Evction's version change
	return o.clientset.Policy().Evictions(eviction.Namespace).Evict(eviction)
}

// deleteOrEvictPods deletes or evicts the pods on the api server
func (o *drainOperation) deleteOrEvictPods(pods []v1.Pod) error {
	if len(pods) == 0 {
		return nil
	}

	policyGroupVersion, err := supportEviction(o.clientset)
	if err != nil {
		return err
	}

	getPodFn := func(namespace, name string) (*v1.Pod, error) {
		return o.clientset.Core().Pods(namespace).Get(name, metav1.GetOptions{})
	}

	if len(policyGroupVersion) > 0 {
		return o.evictPods(pods, policyGroupVersion, getPodFn)
	} else {
		return o.deletePods(pods, getPodFn)
	}
}

func supportEviction(clientset *kubernetes.Clientset) (string, error) {
	discoveryClient := clientset.Discovery()
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
		if resource.Name == EvictionSubresource && resource.Kind == EvictionKind {
			return policyGroupVersion, nil
		}
	}
	return "", nil
}

func (o *drainOperation) evictPods(pods []v1.Pod, policyGroupVersion string, getPodFn func(namespace, name string) (*v1.Pod, error)) error {
	doneCh := make(chan bool, len(pods))
	errCh := make(chan error, 1)

	for _, pod := range pods {
		go func(pod v1.Pod, doneCh chan bool, errCh chan error) {
			var err error
			for {
				err = o.evictPod(pod, policyGroupVersion)
				if err == nil {
					break
				} else if apierrors.IsNotFound(err) {
					doneCh <- true
					return
				} else if apierrors.IsTooManyRequests(err) {
					time.Sleep(5 * time.Second)
				} else {
					errCh <- fmt.Errorf("error when evicting pod %q: %v", pod.Name, err)
					return
				}
			}
			podArray := []v1.Pod{pod}
			_, err = o.waitForDelete(podArray, interval, time.Duration(math.MaxInt64), true, getPodFn)
			if err == nil {
				doneCh <- true
			} else {
				errCh <- fmt.Errorf("error when waiting for pod %q terminating: %v", pod.Name, err)
			}
		}(pod, doneCh, errCh)
	}

	doneCount := 0
	globalTimeout := time.Duration(math.MaxInt64)
	for {
		select {
		case err := <-errCh:
			return err
		case <-doneCh:
			doneCount++
			if doneCount == len(pods) {
				return nil
			}
		case <-time.After(globalTimeout):
			return fmt.Errorf("Drain did not complete within %v", globalTimeout)
		}
	}
}

func (o *drainOperation) deletePods(pods []v1.Pod, getPodFn func(namespace, name string) (*v1.Pod, error)) error {
	globalTimeout := time.Duration(math.MaxInt64)

	for _, pod := range pods {
		err := o.deletePod(pod)
		if err != nil && !apierrors.IsNotFound(err) {
			return err
		}
	}
	_, err := o.waitForDelete(pods, interval, globalTimeout, false, getPodFn)
	return err
}

func (o *drainOperation) waitForDelete(pods []v1.Pod, interval, timeout time.Duration, usingEviction bool, getPodFn func(string, string) (*v1.Pod, error)) ([]v1.Pod, error) {
	var verbStr string
	if usingEviction {
		verbStr = "evicted"
	} else {
		verbStr = "deleted"
	}
	err := wait.PollImmediate(interval, timeout, func() (bool, error) {
		pendingPods := []v1.Pod{}
		for i, pod := range pods {
			p, err := getPodFn(pod.Namespace, pod.Name)
			if apierrors.IsNotFound(err) || (p != nil && p.ObjectMeta.UID != pod.ObjectMeta.UID) {
				o.logger.Infof("%s pod successfully %s", pod.Name, verbStr)
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
