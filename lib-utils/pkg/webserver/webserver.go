package webserver

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]map[string]http.HandlerFunc
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]map[string]http.HandlerFunc),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(method, path string, handler http.HandlerFunc) error {
	if _, ok := s.Handlers[path][method]; ok {
		return fmt.Errorf("there is already a handler for method %s on %s", method, path)
	}
	if _, ok := s.Handlers[path]; !ok {
		s.Handlers[path] = make(map[string]http.HandlerFunc)
	}
	s.Handlers[path][method] = handler
	return nil
}

// Start loop through the handlers and add them to the router, register middleware logger and start the server
func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	for path, handlers := range s.Handlers {
		for method, handler := range handlers {
			s.Router.Method(method, path, handler)
		}
	}
	log.Println("Running server on port", s.WebServerPort)
	if err := http.ListenAndServe(s.WebServerPort, s.Router); err != nil {
		log.Fatal(err)
	}
}
