package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

func (s *Server) createRoutes() {
	r := mux.NewRouter()
	r.HandleFunc("/", s.handleIndex())
	r.HandleFunc("/arc/{arc}", s.handleArc())

	s.routes = r
}

func (s *Server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "/arc/intro", http.StatusSeeOther)
	}
}

func (s *Server) handleArc() http.HandlerFunc {
	t, err := template.ParseFiles("./server/templates/template.html")
	if err != nil {
		log.Fatal("failed to parse template.html", err)
	}
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		arc := vars["arc"]
		arcdata, ok := s.store.Arcs[arc]
		if !ok {
			http.Error(w, fmt.Sprintf("arc %s does not exist", arc), http.StatusBadRequest)
			return
		}
		err = t.Execute(w, arcdata)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to create tmpl for %s", arc), http.StatusInternalServerError)
		}
	}
}
