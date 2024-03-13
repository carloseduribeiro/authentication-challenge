package main

import (
	"context"
	"github.com/carloseduribeiro/auth-challenge/customer-debts/configs"
	"github.com/carloseduribeiro/auth-challenge/customer-debts/internal/infra/database"
	"github.com/carloseduribeiro/auth-challenge/customer-debts/internal/infra/http/handlers"
	"github.com/carloseduribeiro/auth-challenge/lib-utils/pkg/config"
	"github.com/carloseduribeiro/auth-challenge/lib-utils/pkg/db"
	"github.com/carloseduribeiro/auth-challenge/lib-utils/pkg/webserver"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting application...")
	cfg, err := config.LoadConfig[configs.Conf](".")
	if err != nil {
		log.Fatal(err)
	}

	db.RunMigrations(cfg.DatabaseURL)

	dbConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	dbPool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()

	webServer := webserver.NewWebServer(cfg.WebServerPort)
	repository := database.NewDebtsRepository(dbPool)
	createDebtHandler := handlers.NewCreateDebtHandler(repository)
	if err = webServer.AddHandler(http.MethodPost, "/customer/debts", createDebtHandler.Handle); err != nil {
		log.Fatal(err)
	}
	getDebtsHandler := handlers.NewGetDebtsHandler(repository)
	if err = webServer.AddHandler(http.MethodGet, "/customer/debts/{userDocument}", getDebtsHandler.Handle); err != nil {
		log.Fatal(err)
	}

	webServer.Start()
}
