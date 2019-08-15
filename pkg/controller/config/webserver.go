package config

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/gobuffalo/packr"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
)

type linkItem struct {
	URL         string
	Name        string
	Image       string
	Description string
	DocURL      string
}

type webServer struct {
	sync.RWMutex
	links map[string]linkItem
}

func createWebServer() *webServer {
	ws := &webServer{}
	ws.links = make(map[string]linkItem)

	box := packr.NewBox("./templates")
	stringTemplate, err := box.FindString("index.html")
	if err != nil {
		panic(err)
	}
	tpl := template.Must(template.New("index").Parse(stringTemplate))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, ws.links)
	})

	go http.ListenAndServe(":9000", nil)
	return ws
}

func (s *webServer) AddIngress(name string, item *extensionsv1beta1.Ingress) {
	link := ingressToLink(item)
	if link.Name == "" {
		log.Printf("ingress skipped: %s", name)
		return
	}
	log.Printf("ingress added: %s", name)
	s.Lock()
	s.links[name] = link
	s.Unlock()
}

func (s *webServer) Remove(name string) {
	log.Printf("ingress removed: %s", name)
	s.Lock()
	delete(s.links, name)
	s.Unlock()
}

func ingressToLink(ingress *extensionsv1beta1.Ingress) linkItem {
	link := linkItem{
		Name:        ingress.Annotations["kube-linker/name"],
		Description: ingress.Annotations["kube-linker/description"],
		Image:       ingress.Annotations["kube-linker/image"],
		DocURL:      ingress.Annotations["kube-linker/doc-url"],
	}

	if len(ingress.Spec.TLS) == 0 {
		link.URL = fmt.Sprintf("http://%s", ingress.Spec.Rules[0].Host)
	} else {
		link.URL = fmt.Sprintf("https://%s", ingress.Spec.Rules[0].Host)
	}

	return link
}
