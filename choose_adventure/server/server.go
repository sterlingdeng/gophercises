package server

import (
	"fmt"
	"github.com/sterlingdeng/cyoa/story"
	"log"
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	store  story.Story
	routes *mux.Router
	port   int
}

func NewServer(port int, store story.Story) Server {
	return Server{
		store:  store,
		port:   port,
	}
}

func (s *Server) Start() {
	s.createRoutes()
	fmt.Printf("starting server on port %d\n", s.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.routes))
}

