package e2e

import (
	goctx "context"
	"fmt"
	"testing"

	corev1 "k8s.io/api/core/v1"

	framework "github.com/operator-framework/operator-sdk/pkg/test"

	"github.com/DataDog/k8s-dns-exposer/test/e2e/utils"
)

// TestK8SDNSExposerController runs the test suite
func TestController(t *testing.T) {
	// run subtests
	t.Run("controller-group", func(t *testing.T) {
		t.Run("Create-Endpoint", CreateEpFromService)
	})
}

// CreateEndpointFromService test if en Enpoint is properly created from a service
func CreateEpFromService(t *testing.T) {
	t.Parallel()
	f, ctx, err := InitController(t)
	defer ctx.Cleanup()
	if err != nil {

		t.Fatal(err)
	}

	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Fatal(fmt.Errorf("could not get namespace: %v", err))
	}

	serviceName := "foo"
	externalName := "datadoghq.com"

	newService := utils.NewService(namespace, serviceName, externalName, map[string]string{"datadoghq.com/k8s-dns-exposer": "true"})
	err = f.Client.Create(goctx.TODO(), newService, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		t.Fatal(err)
	}

	// test if the canary deployment have been updated
	isUpdated := func(ep *corev1.Endpoints) (bool, error) {
		if ep == nil {
			return false, nil
		}
		for _, subset := range ep.Subsets {
			if len(subset.Addresses) > 0 {
				return true, nil
			}

		}
		return false, nil
	}
	// check the update on the master deployment
	err = utils.WaitForFuncOnEndpoints(t, f.KubeClient, namespace, serviceName, isUpdated, retryInterval, timeout)
	if err != nil {
		t.Fatal(err)
	}
}
