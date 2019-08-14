package config

import (
	"fmt"
	"net/http"
	"sync"

	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
)

type linkItem struct {
	url  string
	desc string
}

// links:  make(map[string]LinkItem),

type webServer struct {
	sync.RWMutex
	links map[string]linkItem
}

func createWebServer() *webServer {
	ws := &webServer{}
	ws.links = make(map[string]linkItem)

	http.HandleFunc("/", ws.handler)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		panic(err)
	}
	return ws
}

func (s *webServer) handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Maksim!")
}

func (s *webServer) AddIngress(name string, item *extensionsv1beta1.Ingress) {
	s.Lock()
	s.links[name] = linkItem{
		url:  "http://www.ya.ru",
		desc: "some description here",
	}
	s.Unlock()
}

func (s *webServer) Remove(name string) {
	s.Lock()
	delete(s.links, name)
	s.Unlock()
}
