package grpc_server

import (
	"fmt"
	"net"

	env_interface "github.com/jialunzhai/crimemap/analytics/online/server/enviroment"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	env      env_interface.Env
	listener net.Listener
	server   *grpc.Server
}

func Register(env env_interface.Env) error {
	s, err := NewGRPCServer(env)
	if err != nil {
		return err
	}
	env.SetGRPCServer(s)
	return nil
}

func NewGRPCServer(env env_interface.Env) (*GRPCServer, error) {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return nil, err
	}
	server := grpc.NewServer()
	return &GRPCServer{
		env:      env,
		listener: listener,
		server:   server,
	}, nil
}

func (s *GRPCServer) GetServer() *grpc.Server {
	return s.server
}

func (s *GRPCServer) Run() error {
	if s.env.GetCrimeMapService() == nil {
		fmt.Errorf("CrimeMapServer.Register should be called before GRPCServer.Run")
	}
	return s.server.Serve(s.listener)
}
