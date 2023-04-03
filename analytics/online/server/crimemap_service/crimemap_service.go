package crimemap_server

import (
	"context"
	"fmt"

	pb "cs.nyu.edu/crimemap/analytics/online/proto"
	env_interface "cs.nyu.edu/crimemap/analytics/online/server/enviroment"
)

type CrimeMapService struct {
	env env_interface.Env
}

func Register(env env_interface.Env) error {
	grpcServer := env.GetGRPCServer()
	if grpcServer == nil {
		return fmt.Errorf("GRPCServer.Register should be called before CrimeMapService.Register")
	}
	s, err := NewCrimeMapService(env)
	if err != nil {
		return err
	}
	pb.RegisterCrimeMapServer(grpcServer.GetServer(), s)
	env.SetCrimeMapService(s)
	return nil
}

func NewCrimeMapService(env env_interface.Env) (*CrimeMapService, error) {
	return &CrimeMapService{
		env: env,
	}, nil
}

func (s *CrimeMapService) GetCrimes(ctx context.Context, _ *pb.GetCrimesRequest) (*pb.GetCrimesResponse, error) {
	// TODO: implement this with trino_client.GetCrimes
	return &pb.GetCrimesResponse{}, nil
}
