package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

// Start will be executed when operator sdk manager starts
// when the channel closes - the whole operator will stop running
func (s *Server) Start(channel <-chan struct{}) error {

	s.links = make(map[string]LinkItem)
	s.template = template.Must(template.New("index").Parse(htmlTemplate))

	r := mux.NewRouter()
	r.HandleFunc("/", s.serverDocHandler).Methods("GET")
	r.HandleFunc("/api/link/{id}", s.apiCreateLinkHandler).Methods("POST")
	r.HandleFunc("/api/link/{id}", s.apiDeleteLinkHandler).Methods("DELETE")

	http.Handle("/", r)
	return http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func (s *Server) apiCreateLinkHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var link LinkItem
	_ = json.NewDecoder(r.Body).Decode(&link)

	s.Lock()
	s.links[params["id"]] = link
	s.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(link)
}

func (s *Server) apiDeleteLinkHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	s.Lock()
	delete(s.links, params["id"])
	s.Unlock()
}

func (s *Server) serverDocHandler(w http.ResponseWriter, r *http.Request) {
	s.template.Execute(w, s.links)
}
