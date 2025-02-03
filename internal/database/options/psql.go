package options

import (
	"context"
	"fmt"
	"log"
	"server/internal/database/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Psql struct{}

var ConnectionPool *pgxpool.Pool

// Initialize the connection pool
func InitPool(ctx context.Context) error {

	var err error
	ConnectionPool, err = pgxpool.NewWithConfig(context.Background(), config.Psql())
	if err != nil {
		return fmt.Errorf("error creating connection pool: %w", err)
	}
	log.Println("Connection pool created")
	return nil

}

func GetConnection(ctx context.Context) (*pgxpool.Conn, error) {

	conn, err := ConnectionPool.Acquire(ctx)

	if err != nil {
		log.Println("error in extracting the connection from pool")
		return nil, fmt.Errorf("database connection failed !!")
	}

	return conn, nil
}

func ClosePool() {
	ConnectionPool.Close()
}
