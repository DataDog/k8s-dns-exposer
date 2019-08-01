package service

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/DataDog/k8s-dns-exposer/pkg/controller/config"
	"github.com/DataDog/k8s-dns-exposer/pkg/controller/predicate"
	"github.com/DataDog/k8s-dns-exposer/pkg/utils"
)

var log = logf.Log.WithName("k8s-dns-exposer")

// Add creates a new Service Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileService{
		client: mgr.GetClient(),
		scheme: mgr.GetScheme(),

		dnsResolver: utils.NewDNSResolver(),
		updateEndpointsFunc: utils.UpdateEndpoints,

		watcherPredicate: predicate.AnnotationPredicate{Key: config.K8sDNSExposerAnnotationKey},
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("k8s-dns-exposer", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Service
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForObject{}, r.(*ReconcileService).watcherPredicate)
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileService implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileService{}

// ReconcileService reconciles a Service object
type ReconcileService struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme

	dnsResolver        utils.DNSResolverIface
	updateEndpointsFunc utils.UpdateEndpointsFunc

	watcherPredicate predicate.AnnotationPredicate
}

// Reconcile reads that state of the cluster for a Service object and makes changes based on the state read
// and what is in the Service.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileService) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Service")

	// use defaultResult in order to force the Service reconcile requeue after the default Requeue duration
	defaultResult := reconcile.Result{
		RequeueAfter: config.DefaultRequeueDuration,
	}

	// Fetch the Service instance
	service := &corev1.Service{}
	if err := r.client.Get(context.TODO(), request.NamespacedName, service); err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// retrieve associated endpoint if exist
	endpoint := &corev1.Endpoints{}
	isEndpointCreationMode := false
	if err := r.client.Get(context.TODO(), request.NamespacedName, endpoint); err != nil {
		if errors.IsNotFound(err) {
			isEndpointCreationMode = true
		} else {
			return defaultResult, err
		}
	}

	ips, err := r.dnsResolver.Resolve(service.Spec.ExternalName)
	if err != nil {
		return defaultResult, err
	}

	var newEndpoint *corev1.Endpoints
	newEndpoint, err = r.updateEndpointsFunc(endpoint, service, r.scheme, ips)
	if err != nil {
		return defaultResult, err
	}

	if isEndpointCreationMode {
		err = r.client.Create(context.TODO(), newEndpoint)
		reqLogger.Info("Endpoint created")
	} else if !apiequality.Semantic.DeepEqual(endpoint, newEndpoint) {
		err = r.client.Update(context.TODO(), newEndpoint)
		reqLogger.Info("Endpoint updated")
	}

	return defaultResult, err
}
