package virtualservice

import (
	"context"

	virtualservice "github.com/afoninsky/kube-linker/pkg/apis/networking/v1alpha3"

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
	c, err := controller.New("virtualservice-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// watch for virtualservices events
	return c.Watch(&source.Kind{Type: &virtualservice.VirtualService{}}, &handler.EnqueueRequestForObject{})

}

// blank assignment to verify that ReconcileConfig implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileConfig{}

// Reconcile handles ingress updates
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileConfig) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	virtualservice := &virtualservice.VirtualService{}
	apiError := r.client.Get(context.TODO(), request.NamespacedName, virtualservice)
	link := virtualserviceToLink(virtualservice)

	if apiError != nil {
		if errors.IsNotFound(apiError) {
			webError := r.web.Do("DELETE", link)
			return reconcile.Result{}, webError
		}
		// error reading the object - requeue the request.
		return reconcile.Result{}, apiError
	}
	// ingress is either created or updated
	// ensure its allowed to display
	_, enabled := virtualservice.Annotations["kube-linker/enabled"]
	if !enabled {
		return reconcile.Result{}, nil
	}

	webError := r.web.Do("POST", link)
	return reconcile.Result{}, webError
}

func virtualserviceToLink(service *virtualservice.VirtualService) web.LinkItem {
	link := web.LinkItem{
		AnnotatedName:        service.Annotations["kube-linker/name"],
		AnnotatedDescription: service.Annotations["kube-linker/description"],
		AnnotatedURL:         service.Annotations["kube-linker/doc-url"],
		SpecName:             service.Name,
		SpecNamespace:        service.Namespace,
		SpecType:             "virtualservice",
		SpecURL:              service.Spec.Hosts,
	}
	// if len(ingress.Spec.TLS) == 0 {
	// 	link.SpecURL = append(link.SpecURL, fmt.Sprintf("http://%s", ingress.Spec.Rules[0].Host))
	// } else {
	// 	link.SpecURL = append(link.SpecURL, fmt.Sprintf("https://%s", ingress.Spec.Rules[0].Host))
	// }

	return link
}
