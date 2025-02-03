package config

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const defaultMaxConns = int32(10)
const defaultMinConns = int32(0)
const defaultMaxConnLifetime = time.Minute * 10
const defaultMaxConnIdleTime = time.Minute * 5
const defaultHealthCheckPeriod = time.Minute
const defaultConnectTimeout = time.Second * 8
const DATABASE_URL string = "postgres://postgres:adityaisprono1@localhost:5432/codespace"

func Psql() *pgxpool.Config {

	dbConfig, err := pgxpool.ParseConfig(DATABASE_URL)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	dbConfig.AfterConnect = func(ctx context.Context, c *pgx.Conn) error {
		return nil
	}

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		log.Println("connection is aquired !!")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		log.Println("connection is released !!")
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		log.Println("Disconnected from database")
	}

	return dbConfig
}
