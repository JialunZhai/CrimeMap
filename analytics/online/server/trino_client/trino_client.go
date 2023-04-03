package trino_client

import (
	"context"
	"database/sql"

	env_interface "github.com/jialunzhai/crimemap/analytics/online/server/enviroment"
	"github.com/jialunzhai/crimemap/analytics/online/server/interfaces"
	_ "github.com/trinodb/trino-go-client/trino"
)

type TrinoClient struct {
	env env_interface.Env
	db  *sql.DB
}

func Register(env env_interface.Env) error {
	s, err := NewTrinoClient(env)
	if err != nil {
		return err
	}
	env.SetTrinoClient(s)
	return nil
}

func NewTrinoClient(env env_interface.Env) (*TrinoClient, error) {
	dsn := "http://user@localhost:8080?catalog=default&schema=test"
	db, err := sql.Open("trino", dsn)
	if err != nil {
		return nil, err
	}
	return &TrinoClient{
		env: env,
		db:  db,
	}, nil
}

func (c *TrinoClient) GetCrimes(ctx context.Context, minX float64, maxX float64, minY float64, maxY float64) ([]*interfaces.Crime, error) {
	crimes := make([]*interfaces.Crime, 0)
	// TODO: fill the SQL statement in QueryContext
	rows, err := c.db.QueryContext(ctx, "SELECT name FROM users WHERE age = $1", 100)
	if err != nil {
		return crimes, err
	}
	defer rows.Close()

	for rows.Next() {
		var crime interfaces.Crime
		if err := rows.Scan(&crime); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			return crimes, err
		}
		crimes = append(crimes, &crime)
	}
	// If the database is being written to ensure to check for Close
	// errors that may be returned from the driver. The query may
	// encounter an auto-commit error and be forced to rollback changes.
	rerr := rows.Close()
	if rerr != nil {
		return crimes, err
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		return crimes, err
	}
	return crimes, nil
}

func (c *TrinoClient) Close() error {
	return c.db.Close()
}
