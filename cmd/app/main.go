package main

import (
	"log"
	"os"

	"github.com/dannamer/JavaCode-test/internal/api"
	"github.com/dannamer/JavaCode-test/internal/repository/postgresql"
	"github.com/dannamer/JavaCode-test/internal/service"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	config, err := postgresql.NewConfig()
	if err != nil {
		log.Fatal("Ошибка конфигурации PostgreSQL:", err)
	}

	runDBMigration(os.Getenv("MIGRATION_URL"), config.GetDSN())

	postgres, err := postgresql.NewPostgres(*config)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	repo := postgresql.NewWalletRepo(postgres.Pool)
	serv := service.NewWalletService(&repo)
	server := api.NewWalletHandler(&serv)

	server.RunServer()
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create a new migrate instance", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up:", err)
	}
	log.Println("db migrated successfully")
}
