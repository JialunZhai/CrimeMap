package http_server

import (
	"net/http"

	env_interface "github.com/jialunzhai/crimemap/analytics/online/server/enviroment"
)

type HTTPServer struct {
	env env_interface.Env
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
	return &HTTPServer{
		env: env,
	}, nil
}

func (s *HTTPServer) Run() error {
	http.Handle("/", http.FileServer(http.Dir("./analytics/online/webapp/app")))
	return http.ListenAndServe(":8080", nil)
}
