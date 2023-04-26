package hbase_client

import (
	"context"
	"errors"
	"log"

	env_interface "github.com/jialunzhai/crimemap/analytics/online/server/enviroment"
	"github.com/jialunzhai/crimemap/analytics/online/server/interfaces"
	"github.com/tsuna/gohbase"
)

type HBaseClient struct {
	env    env_interface.Env
	client gohbase.Client
}

func Register(env env_interface.Env) error {
	config := env.GetConfig()
	if config == nil || config.Database.Address == "" {
		return errors.New("HBase client not configured")
	}
	s, err := NewHBaseClient(env, config.Database.Address)
	if err != nil {
		return err
	}
	env.SetDatabaseClient(s)
	return nil
}

func NewHBaseClient(env env_interface.Env, zkquorum string) (*HBaseClient, error) {
	client := gohbase.NewClient(zkquorum)
	log.Println("Connected to HBase server.")
	return &HBaseClient{
		env:    env,
		client: client,
	}, nil
}

func (c *HBaseClient) GetCrimes(ctx context.Context, minX float64, maxX float64, minY float64, maxY float64) ([]*interfaces.Crime, error) {
	crimes := make([]*interfaces.Crime, 0)
	return crimes, nil
}

func (c *HBaseClient) Close() error {
	c.client.Close()
	return nil
}
