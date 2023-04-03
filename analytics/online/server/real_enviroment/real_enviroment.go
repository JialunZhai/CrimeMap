package real_environment

import (
	"github.com/jialunzhai/crimemap/analytics/online/server/interfaces"
)

type RealEnv struct {
	httpServer      interfaces.HTTPServer
	grpcServer      interfaces.GRPCServer
	crimemapService interfaces.CrimeMapService
	trinoClient     interfaces.TrinoClient
}

func NewRealEnv() *RealEnv {
	return &RealEnv{}
}

func (r *RealEnv) GetHTTPServer() interfaces.HTTPServer {
	return r.httpServer
}

func (r *RealEnv) SetHTTPServer(s interfaces.HTTPServer) {
	r.httpServer = s
}

func (r *RealEnv) GetGRPCServer() interfaces.GRPCServer {
	return r.grpcServer
}

func (r *RealEnv) SetGRPCServer(s interfaces.GRPCServer) {
	r.grpcServer = s
}

func (r *RealEnv) GetCrimeMapService() interfaces.CrimeMapService {
	return r.crimemapService
}

func (r *RealEnv) SetCrimeMapService(s interfaces.CrimeMapService) {
	r.crimemapService = s
}

func (r *RealEnv) GetTrinoClient() interfaces.TrinoClient {
	return r.trinoClient
}

func (r *RealEnv) SetTrinoClient(s interfaces.TrinoClient) {
	r.trinoClient = s
}
