package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func ConnDB() (*pgx.Conn, error) {
	var err error

	if err = godotenv.Load(".env"); err != nil {
		return nil, err
	}

	var (
		databaseUser string = os.Getenv("DATABASE_USER")
		databasePass string = os.Getenv("DATABASE_PASS")
		databasePort string = os.Getenv("DATABASE_PORT")
		databaseName string = os.Getenv("DATABASE_NAME")
	)

	connStr := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", databaseUser, databasePass, databasePort, databaseName)

	conn, err := pgx.Connect(context.Background(), connStr)

	if err != nil {
		return nil, err
	}

	defer conn.Close(context.Background())

	return conn, nil
}
