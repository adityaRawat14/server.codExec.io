package db

import (
	"context"
	"log"
	"server/internal/database/options"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New() (*pgxpool.Conn, error) {

	c, err := options.GetConnection(context.Background()) 

	if err != nil {
		log.Println("error : New() -> GetConnection():", err)
		return nil, err
	}

	return c, nil
}

func Close(conn *pgxpool.Conn) {
	conn.Release()
}

func OpenDbPool() error {
	err := options.InitPool(context.Background())
	if err != nil {
		log.Println("error : OpenDbPool() :", err)
		return err
	}
	return nil
}

func ShutDownDbPool() {
	options.ClosePool()
}
