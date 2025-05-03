package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/dafraer/effective-mobile-task/api"
	"github.com/dafraer/effective-mobile-task/enrich"
	"github.com/dafraer/effective-mobile-task/store"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	//Create logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Errorf("error while creating new Logger, %v ", err))
	}
	sugar := logger.Sugar()
	defer func() {
		if err := sugar.Sync(); err != nil {
			panic(err)
		}
	}()

	//Load environment variables
	err = godotenv.Load()
	if err != nil {
		panic(err)
	}
	address := os.Getenv("ADDRESS")
	dbConnStr := os.Getenv("DB_URI")
	if address == "" || dbConnStr == "" {
		panic("error loading environment variables")
	}
	sugar.Infow("Loaded environment variables")

	//Connect to the db and perform migrations
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		panic(err)
	}
	sugar.Infow("Connection to the database established")

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "postgres", driver)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}
	storage := store.New(db, sugar)
	sugar.Infow("Migrations performed")

	//Create enricher
	enricher := enrich.New(sugar)

	//Create and run the service
	service := api.New(sugar, storage, enricher)
	sugar.Infow("New service created")

	if err := service.Run(context.Background(), address); err != nil {
		panic(err)
	}
	sugar.Infow("Service stopped")
}
