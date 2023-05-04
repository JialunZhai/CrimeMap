package prometheus_monitor

import (
	"context"
	"errors"
	"log"
	"net/http"

	env_interface "github.com/jialunzhai/crimemap/analytics/online/server/enviroment"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusMonitor struct {
	env    env_interface.Env
	server *http.Server
}

var ErrorCounter = promauto.NewCounter(prometheus.CounterOpts{
	Name: "error_counter",
	Help: "The total number of errors occurred.",
})

func Register(env env_interface.Env) error {
	config := env.GetConfig()
	if config == nil || config.Prometheus.Address == "" {
		return errors.New("Prometheus monitor not configured")
	}
	m, err := NewPrometheusMonitor(env, config.Prometheus.Address)
	if err != nil {
		return err
	}
	env.SetPrometheusMonitor(m)
	return nil
}

func NewPrometheusMonitor(env env_interface.Env, address string) (*PrometheusMonitor, error) {
	server := &http.Server{
		Addr:    address,
		Handler: promhttp.Handler(),
	}
	return &PrometheusMonitor{
		env:    env,
		server: server,
	}, nil
}

func (m *PrometheusMonitor) Run() error {
	log.Printf("Start to serve Prometheus scrapes from address: %v\n", m.server.Addr)
	return m.server.ListenAndServe()
}

func (m *PrometheusMonitor) Shutdown(ctx context.Context) error {
	return m.server.Shutdown(ctx)
}
