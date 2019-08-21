package webserver

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"

	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
)

type linkItem struct {
	// kube-linker/name
	AnnotatedName string
	// kube-linker/description
	AnnotatedDescription string
	// kube-linker/doc-url
	AnnotatedURL  string
	SpecURL       []string
	SpecName      string
	SpecNamespace string
	SpecEndpoint  string
	SpecType      string
}

type WebServer struct {
	sync.RWMutex
	links map[string]linkItem
}

func New() *WebServer {
	ws := &WebServer{}
	ws.links = make(map[string]linkItem)
	tpl := template.Must(template.New("index").Parse(htmlTemplate))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, ws.links)
	})

	go http.ListenAndServe(":9000", nil)
	return ws
}

func (s *WebServer) AddIngress(name string, item *extensionsv1beta1.Ingress) {
	link := ingressToLink(item)
	_, enabled := item.Annotations["kube-linker/enabled"]
	if !enabled {
		log.Printf("ingress skipped: %s", name)
		return
	}
	log.Printf("ingress added: %s", name)
	s.Lock()
	s.links[name] = link
	s.Unlock()
}

func (s *WebServer) Remove(name string) {
	log.Printf("ingress removed: %s", name)
	s.Lock()
	delete(s.links, name)
	s.Unlock()
}

func ingressToLink(ingress *extensionsv1beta1.Ingress) linkItem {
	link := linkItem{
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
