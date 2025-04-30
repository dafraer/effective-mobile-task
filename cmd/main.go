package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/dafraer/effective-mobile-task/api"
	"github.com/dafraer/effective-mobile-task/store"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	//Load environment veriables
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	address := os.Getenv("ADDRESS")
	dbConnStr := os.Getenv("DB_URI")
	if address == "" || dbConnStr == "" {
		panic("error loading environment variables")
	}

	//Create logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Errorf("error while creating new Logger, %v ", err))
	}
	sugar := logger.Sugar()
	defer sugar.Sync()

	//Connect to the db nad perfprm migrations
	db, err := sql.Open("postgres", dbConnStr)
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver)
	m.Up()
	storage := store.New(db)

	//Create and run the servie
	service := api.New(sugar, storage)
	if err := service.Run(context.Background(), address); err != nil {
		panic(err)
	}
}
