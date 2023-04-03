package main

import (
	"context"
	"log"

	cms "github.com/jialunzhai/crimemap/analytics/online/server/crimemap_service"
	"github.com/jialunzhai/crimemap/analytics/online/server/grpc_server"
	"github.com/jialunzhai/crimemap/analytics/online/server/http_server"
	real_env "github.com/jialunzhai/crimemap/analytics/online/server/real_enviroment"
	"github.com/jialunzhai/crimemap/analytics/online/server/trino_client"
	"golang.org/x/sync/errgroup"
)

const debug = true

func main() {
	env := real_env.NewRealEnv()

	if debug {
		if err := http_server.Register(env); err != nil {
			log.Fatalf("HTTPServer.Register failed: %v\n", err)
		}
		if err := env.GetHTTPServer().Run(); err != nil {
			log.Fatalf("HTTPServer exited because %v\n", err)
		}
		return
	}

	if err := trino_client.Register(env); err != nil {
		log.Fatalf("TrinoClient.Register failed: %v\n", err)
	}
	if err := grpc_server.Register(env); err != nil {
		log.Fatalf("GRPCServer.Register failed: %v\n", err)
	}
	if err := cms.Register(env); err != nil {
		log.Fatalf("CrimeMapServer.Register failed: %v\n", err)
	}
	if err := http_server.Register(env); err != nil {
		log.Fatalf("HTTPServer.Register failed: %v\n", err)
	}

	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		err := env.GetGRPCServer().Run()
		log.Fatalf("GRPCServer exited because %v\n", err)
		return err
	})
	g.Go(func() error {
		err := env.GetHTTPServer().Run()
		log.Fatalf("HTTPServer exited because %v\n", err)
		return err
	})
}
