package environment

import (
	"github.com/jialunzhai/crimemap/analytics/online/server/interfaces"
)

type Env interface {
	GetHTTPServer() interfaces.HTTPServer
	SetHTTPServer(interfaces.HTTPServer)
	GetGRPCServer() interfaces.GRPCServer
	SetGRPCServer(interfaces.GRPCServer)
	GetCrimeMapService() interfaces.CrimeMapService
	SetCrimeMapService(interfaces.CrimeMapService)
	GetTrinoClient() interfaces.TrinoClient
	SetTrinoClient(interfaces.TrinoClient)
}
