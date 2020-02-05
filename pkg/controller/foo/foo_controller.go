package foo

import (
	"context"

	samplecontrollerv1alpha1 "github.com/govargo/foo-controller-operatorsdk/pkg/apis/samplecontroller/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_foo")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Foo Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileFoo{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("foo-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Foo
	err = c.Watch(&source.Kind{Type: &samplecontrollerv1alpha1.Foo{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Deployment and requeue the owner Foo
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &samplecontrollerv1alpha1.Foo{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileFoo implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileFoo{}

// ReconcileFoo reconciles a Foo object
type ReconcileFoo struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Foo object and makes changes based on the state read
// and what is in the Foo.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileFoo) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Foo")
	ctx := context.Background()

	/*
		### 1: Load the Foo by name
		We'll fetch the Foo using our client.
		All client methods take a context (to allow for cancellation) as
		their first argument, and the object
		in question as their last.
		Get is a bit special, in that it takes a
		[`NamespacedName`](https://godoc.org/sigs.k8s.io/controller-runtime/pkg/client#ObjectKey)
		as the middle argument (most don't have a middle argument, as we'll see below).
		Many client methods also take variadic options at the end.
	*/
	foo := &samplecontrollerv1alpha1.Foo{}
	if err := r.client.Get(ctx, request.NamespacedName, foo); err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.Info("Foo not found. Ignore not found")
			return reconcile.Result{}, nil
		}
		reqLogger.Error(err, "failed to get Foo")
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	/*
		### 2: Clean Up old Deployment which had been owned by Foo Resource.
		We'll find deployment object which foo object owns.
		If there is a deployment which is owned by foo and it doesn't match foo.spec.deploymentName,
		we clean up the deployment object.
		(If we do nothing without this func, the old deployment object keeps existing.)
	*/
	if err := r.cleanupOwnedResources(ctx, foo); err != nil {
		reqLogger.Error(err, "failed to clean up old Deployment resources for this Foo")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

// cleanupOwnedResources will Delete any existing Deployment resources that
// were created for the given Foo that no longer match the
// foo.spec.deploymentName field.
func (r *ReconcileFoo) cleanupOwnedResources(ctx context.Context, foo *samplecontrollerv1alpha1.Foo) error {
	reqLogger := log.WithValues("Request.Namespace", foo.Namespace, "Request.Name", foo.Name)
	reqLogger.Info("finding existing Deployments for Foo resource")

	// List all deployment resources owned by this Foo
	deployments := &appsv1.DeploymentList{}
	labelSelector := labels.SelectorFromSet(labelsForFoo(foo.Name))
	listOps := &client.ListOptions{
		Namespace:     foo.Namespace,
		LabelSelector: labelSelector,
	}
	if err := r.client.List(ctx, deployments, listOps); err != nil{
		reqLogger.Error(err,"faild to get list of deployments")
		return err
	}

	// Delete deployment if the deployment name doesn't match foo.spec.deploymentName
	for _, deployment := range deployments.Items {
		if deployment.Name == foo.Spec.DeploymentName {
			// If this deployment's name matches the one on the Foo resource
			// then do not delete it.
			continue
		}

		// Delete old deployment object which doesn't match foo.spec.deploymentName
		if err := r.client.Delete(ctx, &deployment); err != nil {
			reqLogger.Error(err, "failed to delete Deployment resource")
			return err
		}

		reqLogger.Info("deleted old Deployment resource for Foo", "deploymentName", deployment.Name)
	}

	return nil
}

// labelsForFoo returns the labels for selecting the resources
// belonging to the given foo CR name.
func labelsForFoo(name string) map[string]string {
	return map[string]string{"app": "nginx", "controller": name}
}