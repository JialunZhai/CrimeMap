package http_server

import (
	"context"
	"log"
	"net/http"

	env_interface "github.com/jialunzhai/crimemap/analytics/online/server/enviroment"
)

type HTTPServer struct {
	env    env_interface.Env
	server *http.Server
}

func Register(env env_interface.Env) error {
	s, err := NewHTTPServer(env)
	if err != nil {
		return err
	}
	env.SetHTTPServer(s)
	return nil
}

func NewHTTPServer(env env_interface.Env) (*HTTPServer, error) {
	server := &http.Server{
		Addr:    "localhost:8081",
		Handler: http.FileServer(http.Dir("./analytics/online/webapp/app")),
	}
	return &HTTPServer{
		env:    env,
		server: server,
	}, nil
}

func (s *HTTPServer) Run() error {
	log.Printf("Start to serve HTTP requests from address: %v\n", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
