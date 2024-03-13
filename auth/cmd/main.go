package main

import (
	"github.com/carloseduribeiro/auth-challenge/auth/configs"
	"github.com/carloseduribeiro/auth-challenge/auth/internal/infra/database"
	"github.com/carloseduribeiro/auth-challenge/auth/internal/infra/webserver"
	"github.com/carloseduribeiro/auth-challenge/auth/internal/infra/webserver/http/handlers"
	"github.com/google/uuid"
	"golang.org/x/net/context"
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

	ctx := context.TODO()
	dbPool := configs.SetupDatabase(ctx, config)
	defer dbPool.Close()

	webServer := webserver.NewWebServer(config.WebServerPort)

	userRepository := database.NewUserRepository(dbPool)
	createUserHandler := handlers.NewCreateUser(userRepository, uuid.NewUUID)
	if err = webServer.AddHandler(http.MethodPost, "/auth/users", createUserHandler.Handler); err != nil {
		log.Fatal(err)
	}

	sessionRepository := database.NewSessionRepository(dbPool)
	loginHandler := handlers.NewLogin(userRepository, sessionRepository, config.JWTSecretKey, config.SessionMaxDuration)
	if err = webServer.AddHandler(http.MethodPost, "/auth/users/login", loginHandler.Handler); err != nil {
		log.Fatal(err)
	}

	webServer.Start()
}
