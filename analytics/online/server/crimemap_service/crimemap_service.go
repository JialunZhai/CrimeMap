package crimemap_server

import (
	"context"
	"fmt"

	cmspb "github.com/jialunzhai/crimemap/analytics/online/proto/crimemap_service"
	env_interface "github.com/jialunzhai/crimemap/analytics/online/server/enviroment"
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
	cmspb.RegisterCrimeMapServer(grpcServer.GetServer(), s)
	env.SetCrimeMapService(s)
	return nil
}

func NewCrimeMapService(env env_interface.Env) (*CrimeMapService, error) {
	return &CrimeMapService{
		env: env,
	}, nil
}

func (s *CrimeMapService) GetCrimes(ctx context.Context, _ *cmspb.GetCrimesRequest) (*cmspb.GetCrimesResponse, error) {
	// TODO: implement this with env.GetDatabaseClient().GetCrimes
	return &cmspb.GetCrimesResponse{}, nil
}
