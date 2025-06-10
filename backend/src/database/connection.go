package database

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var (
	DB   *pgxpool.Pool
	once sync.Once
)

func InitDB() error {
	var err error

	once.Do(func() {
		err = godotenv.Load(".env")
		if err != nil {
			return
		}

		databaseUser := os.Getenv("DATABASE_USER")
		databasePass := os.Getenv("DATABASE_PASS")
		databasePort := os.Getenv("DATABASE_PORT")
		databaseName := os.Getenv("DATABASE_NAME")

		connStr := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", databaseUser, databasePass, databasePort, databaseName)

		config, configErr := pgxpool.ParseConfig(connStr)
		if configErr != nil {
			err = configErr
			return
		}

		pool, poolErr := pgxpool.New(context.Background(), config.ConnString())
		if poolErr != nil {
			err = poolErr
			return
		}

		DB = pool
	})

	return err
}
