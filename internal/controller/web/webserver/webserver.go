package webserver

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type HandlerInfo struct {
	Method  string
	Handler http.HandlerFunc
}

type WebServer struct {
	Router        *chi.Mux
	Handlers      map[string][]HandlerInfo
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string][]HandlerInfo),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(method string, path string, handler http.HandlerFunc) {
	s.Handlers[path] = append(s.Handlers[path], HandlerInfo{Method: method, Handler: handler})
}

func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)

	for path, handlers := range s.Handlers {
		for _, info := range handlers {
			switch info.Method {
			case http.MethodGet:
				s.Router.Get(path, info.Handler)
			case http.MethodPost:
				s.Router.Post(path, info.Handler)
			case http.MethodPut:
				s.Router.Put(path, info.Handler)
			case http.MethodDelete:
				s.Router.Delete(path, info.Handler)
			default:
				s.Router.MethodFunc(info.Method, path, info.Handler)
			}
		}
	}

	if err := http.ListenAndServe(s.WebServerPort, s.Router); err != nil {
		log.Fatal(err)
	}
}
