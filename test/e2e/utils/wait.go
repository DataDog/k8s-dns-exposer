package utils

import (
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

//EndpointCheckFunc function to perform checks on endpoint object
type EndpointCheckFunc func(*corev1.Endpoints) (bool, error)

// WaitForFuncOnEndpoints used to wait a valid condition on Endpoints
func WaitForFuncOnEndpoints(t *testing.T, kubeclient kubernetes.Interface, namespace, name string, f EndpointCheckFunc, retryInterval, timeout time.Duration) error {
	return wait.Poll(retryInterval, timeout, func() (bool, error) {
		eps, err := kubeclient.CoreV1().Endpoints(namespace).Get(name, metav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				t.Logf("Waiting for availability of %s endpoint\n", name)
				return false, nil
			}
			return false, err
		}

		ok, err := f(eps)
		t.Logf("Waiting for condition function to be true ok for %s endpoint (%t/%v)\n", name, ok, err)
		return ok, err
	})
}
