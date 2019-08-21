package config

import (
	"log"

	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"

	"k8s.io/apimachinery/pkg/runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	virtualservice "github.com/afoninsky/kube-linker/pkg/apis/networking/v1alpha3"
	ws "github.com/afoninsky/kube-linker/pkg/webserver"
)

// ReconcileConfig reconciles a Config object
type ReconcileConfig struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
	server *ws.WebServer
}

// LinkItem describes link fetched from ingress / virtualservice / etc ...

// Add creates a new Config Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileConfig{
		client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
		server: ws.New(),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New("config-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// watch for ingresses events
	if err := c.Watch(&source.Kind{Type: &extensionsv1beta1.Ingress{}}, &handler.EnqueueRequestForObject{}); err != nil {
		return err
	}

	// watch for ingresses events
	if err := c.Watch(&source.Kind{Type: &virtualservice.VirtualService{}}, &handler.EnqueueRequestForObject{}); err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileConfig implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileConfig{}

// Reconcile handles ingress updates
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileConfig) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log.Print(request)
	// handle ingress event
	// instance := &extensionsv1beta1.Ingress{}
	// err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	// if err != nil {
	// 	if errors.IsNotFound(err) {
	// 		r.server.Remove(request.NamespacedName.String())
	// 		return reconcile.Result{}, nil
	// 	}
	// 	// error reading the object - requeue the request.
	// 	return reconcile.Result{}, err
	// }
	// // ingress is either created or updated
	// r.server.AddIngress(request.NamespacedName.String(), instance)

	return reconcile.Result{}, nil
}