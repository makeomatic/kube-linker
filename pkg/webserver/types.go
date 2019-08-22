package web

import (
	"html/template"
	"sync"
)

type LinkItem struct {
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

type Server struct {
	sync.RWMutex
	links    map[string]LinkItem
	template *template.Template
}

type Client struct{}
