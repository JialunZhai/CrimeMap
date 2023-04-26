package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	cms "github.com/jialunzhai/crimemap/analytics/online/server/crimemap_service"
	"github.com/jialunzhai/crimemap/analytics/online/server/grpc_server"
	"github.com/jialunzhai/crimemap/analytics/online/server/http_server"
	real_env "github.com/jialunzhai/crimemap/analytics/online/server/real_enviroment"
	"github.com/jialunzhai/crimemap/analytics/online/server/hbase_client"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	env := real_env.NewRealEnv()
	if err := hbase_client.Register(env); err != nil {
		log.Fatalf("HBaseClient.Register failed with error: `%v`\n", err)
	}
	if err := grpc_server.Register(env); err != nil {
		log.Fatalf("GRPCServer.Register failed with error: `%v`\n", err)
	}
	if err := cms.Register(env); err != nil {
		log.Fatalf("CrimeMapServer.Register failed with error: `%v`\n", err)
	}
	if err := http_server.Register(env); err != nil {
		log.Fatalf("HTTPServer.Register failed with error: `%v`\n", err)
	}

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		err := env.GetGRPCServer().Run()
		if err != nil {
			log.Printf("GRPCServer shutdowned with error: `%v`\n", err)
		}
		log.Printf("GRPCServer gracefully shutdowned\n")
		return err
	})
	g.Go(func() error {
		err := env.GetHTTPServer().Run()
		if err != http.ErrServerClosed {
			log.Printf("HTTPServer shutdowned with error: `%v`\n", err)
			return err
		}
		log.Printf("HTTPServer gracefully shutdowned\n")
		return err
	})

	// wait for signals
	select {
	case sig := <-sigs:
		// received signal, cancel context
		log.Printf("Received signal: `%v`\n", sig)
		env.GetHTTPServer().Shutdown(ctx)
		env.GetGRPCServer().Shutdown()
		cancel()
		break
	case <-ctx.Done():
		// context cancelled, all goroutines have returned
		break
	}

	if err := env.GetDatabaseClient().Close(); err != nil {
		log.Fatalf("DatabaseClient closed with error: `%v`\n", err)
	}
	log.Printf("DatabaseClient gracefully closed\n")

	// wait for all go-routines in errgroup to return
	if err := g.Wait(); err != nil {
		log.Printf("main exited with error: `%v`\n", err)
	}
}
