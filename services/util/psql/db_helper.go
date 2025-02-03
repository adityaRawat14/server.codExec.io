package psql

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateUser Table creates the users table in the database
func CreateUserTable(db *pgxpool.Conn) error {

	var err error

	// first create extension for the table to provide some randon ids

	err = CreateExtension(db)
	if err != nil {
		log.Println("error : CreateExtension():", err)
		return err
	}

	createTableSQL := `
    CREATE TABLE IF NOT EXISTS users (
   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    image TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
    );`

	_, err = db.Exec(context.Background(), createTableSQL)
	if err != nil {
		fmt.Println("error:CreateUserTable():", err)
		return err
	}

	fmt.Println("table user created sucessfully !!")
	return nil
}

func CreateExtension(db *pgxpool.Conn) error {

	var err error

	createTableSQL := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`

	_, err = db.Exec(context.Background(), createTableSQL)
	if err != nil {
		fmt.Println("error:CreateUserTable():", err)
		return err
	}

	return nil
}
