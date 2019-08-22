package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
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
	r.HandleFunc("/", s.apiCreateLinkHandler).Methods("POST")
	r.HandleFunc("/", s.apiDeleteLinkHandler).Methods("DELETE")

	return http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}

func (s *Server) apiCreateLinkHandler(w http.ResponseWriter, r *http.Request) {
	var link LinkItem
	_ = json.NewDecoder(r.Body).Decode(&link)
	id := getID(link)
	log.Println("item updated:", id)
	s.Lock()
	s.links[id] = link
	s.Unlock()
}

func (s *Server) apiDeleteLinkHandler(w http.ResponseWriter, r *http.Request) {
	var link LinkItem
	_ = json.NewDecoder(r.Body).Decode(&link)
	id := getID(link)
	log.Println("item deleted:", id)

	s.Lock()
	delete(s.links, id)
	s.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&LinkItem{})

}

func (s *Server) serverDocHandler(w http.ResponseWriter, r *http.Request) {
	s.template.Execute(w, s.links)
}

func getID(item LinkItem) string {
	return fmt.Sprintf("%s/%s/%s", item.SpecType, item.SpecNamespace, item.SpecName)
}
