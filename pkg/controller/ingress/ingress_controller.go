package config

import (
	"context"
	"fmt"

	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	web "github.com/afoninsky/kube-linker/pkg/webserver"
)

// ReconcileConfig reconciles a Config object
type ReconcileConfig struct {
	client client.Client
	scheme *runtime.Scheme
	web    web.Client
}

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
		web:    web.NewClient(),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New("ingress-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// watch for ingresses events
	return c.Watch(&source.Kind{Type: &extensionsv1beta1.Ingress{}}, &handler.EnqueueRequestForObject{})

}

// blank assignment to verify that ReconcileConfig implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileConfig{}

// Reconcile handles ingress updates
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileConfig) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	instance := &extensionsv1beta1.Ingress{}
	id := fmt.Sprintf("%s/%s", "ingress", request.NamespacedName.String())
	apiError := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if apiError != nil {
		if errors.IsNotFound(apiError) {
			// TODO: remove item
			webError := r.web.Delete(id)
			return reconcile.Result{}, webError
		}
		// error reading the object - requeue the request.
		return reconcile.Result{}, apiError
	}
	// ingress is either created or updated
	// ensure its allowed to display
	_, enabled := instance.Annotations["kube-linker/enabled"]
	if !enabled {
		return reconcile.Result{}, nil
	}

	webError := r.web.Upsert(id, ingressToLink(instance))
	return reconcile.Result{}, webError
}

func ingressToLink(ingress *extensionsv1beta1.Ingress) web.LinkItem {
	link := web.LinkItem{
		AnnotatedName:        ingress.Annotations["kube-linker/name"],
		AnnotatedDescription: ingress.Annotations["kube-linker/description"],
		AnnotatedURL:         ingress.Annotations["kube-linker/doc-url"],
		SpecName:             ingress.Name,
		SpecNamespace:        ingress.Namespace,
		SpecType:             "ingress",
		SpecEndpoint:         ingress.Spec.Rules[0].HTTP.Paths[0].Backend.ServiceName,
	}
	if len(ingress.Spec.TLS) == 0 {
		link.SpecURL = append(link.SpecURL, fmt.Sprintf("http://%s", ingress.Spec.Rules[0].Host))
	} else {
		link.SpecURL = append(link.SpecURL, fmt.Sprintf("https://%s", ingress.Spec.Rules[0].Host))
	}

	return link
}
