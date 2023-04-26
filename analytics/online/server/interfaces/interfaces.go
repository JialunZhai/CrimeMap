package interfaces

import (
	"context"

	cmspb "github.com/jialunzhai/crimemap/analytics/online/proto/crimemap_service"
	"google.golang.org/grpc"
)

type Crime struct {
	longitude float64
	laitude   float64
	time      int64
}

type HTTPServer interface {
	Run() error
	Shutdown(ctx context.Context) error
}

type GRPCServer interface {
	GetServer() *grpc.Server
	Run() error
	Shutdown()
}

type CrimeMapService interface {
	GetCrimes(ctx context.Context, req *cmspb.GetCrimesRequest) (*cmspb.GetCrimesResponse, error)
}

type DatabaseClient interface {
	GetCrimes(ctx context.Context, minX float64, maxX float64, minY float64, maxY float64) ([]*Crime, error)
	Close() error
}
