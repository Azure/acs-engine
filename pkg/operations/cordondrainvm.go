package operations

import (
	"fmt"
	"math"
	"time"

	"github.com/Azure/acs-engine/pkg/armhelpers"
	log "github.com/sirupsen/logrus"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/pkg/api/v1"
)

const (
	interval            = time.Second * 1
	timeout             = time.Minute * 60
	mirrorPodAnnotation = "kubernetes.io/config.mirror"
)

type drainOperation struct {
	client  armhelpers.KubernetesClient
	node    *v1.Node
	logger  *log.Entry
	timeout time.Duration
}

type podFilter func(v1.Pod) bool

// SafelyDrainNode safely drains a node so that it can be deleted from the cluster
func SafelyDrainNode(az armhelpers.ACSEngineClient, logger *log.Entry, masterURL, kubeConfig, nodeName string) error {
	//get client using kubeconfig
	client, err := az.GetKubernetesClient(masterURL, kubeConfig, interval, timeout)
	if err != nil {
		return err
	}

	//Mark the node unschedulable
	node, err := client.GetNode(nodeName)
	if err != nil {
		return err
	}
	node.Spec.Unschedulable = true
	node, err = client.UpdateNode(node)
	if err != nil {
		return err
	}
	logger.Infof("Node %s has been marked unschedulable.", nodeName)

	//Evict pods in node
	drainOp := &drainOperation{client: client, node: node, logger: logger}
	return drainOp.deleteOrEvictPodsSimple()
}

func (o *drainOperation) deleteOrEvictPodsSimple() error {
	pods, err := o.getPodsForDeletion()
	if err != nil {
		return err
	}
	o.logger.Infof("%d pods need to be removed/deleted", len(pods))

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
	if _, found := pod.ObjectMeta.Annotations[mirrorPodAnnotation]; found {
		return false
	}
	return true
}

// getPodsForDeletion returns all the pods we're going to delete.  If there are
// any pods preventing us from deleting, we return that list in an error.
func (o *drainOperation) getPodsForDeletion() (pods []v1.Pod, err error) {
	podList, err := o.client.ListPods(o.node)
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

// deleteOrEvictPods deletes or evicts the pods on the api server
func (o *drainOperation) deleteOrEvictPods(pods []v1.Pod) error {
	if len(pods) == 0 {
		return nil
	}

	policyGroupVersion, err := o.client.SupportEviction()
	if err != nil {
		return err
	}

	if len(policyGroupVersion) > 0 {
		return o.evictPods(pods, policyGroupVersion)
	}
	return o.deletePods(pods)

}

func (o *drainOperation) evictPods(pods []v1.Pod, policyGroupVersion string) error {
	doneCh := make(chan bool, len(pods))
	errCh := make(chan error, 1)

	for _, pod := range pods {
		go func(pod v1.Pod, doneCh chan bool, errCh chan error) {
			var err error
			for {
				err = o.client.EvictPod(&pod, policyGroupVersion)
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
			_, err = o.client.WaitForDelete(o.logger, podArray, true)
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

func (o *drainOperation) deletePods(pods []v1.Pod) error {
	for _, pod := range pods {
		err := o.client.DeletePod(&pod)
		if err != nil && !apierrors.IsNotFound(err) {
			return err
		}
	}
	_, err := o.client.WaitForDelete(o.logger, pods, false)
	return err
}
