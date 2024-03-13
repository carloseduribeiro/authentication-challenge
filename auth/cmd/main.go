package main

import (
	"github.com/carloseduribeiro/auth-challenge/auth/configs"
	"github.com/carloseduribeiro/auth-challenge/auth/internal/infra/database"
	"github.com/carloseduribeiro/auth-challenge/auth/internal/infra/webserver"
	"github.com/carloseduribeiro/auth-challenge/auth/internal/infra/webserver/http/handlers"
	"github.com/google/uuid"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting application...")
	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	configs.RunMigrations(config.DatabaseURL)

	dbPool := configs.SetupDatabase(config)
	defer dbPool.Close()

	userRepository := database.NewUserRepository(dbPool)
	createUserHandler := handlers.NewCreateUser(userRepository, uuid.NewUUID)
	webServer := webserver.NewWebServer(config.WebServerPort)
	if err = webServer.AddHandler(http.MethodPost, "/auth/users", createUserHandler.Handler); err != nil {
		log.Fatal(err)
	}
	webServer.Start()
}
