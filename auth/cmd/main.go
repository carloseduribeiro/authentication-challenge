package main

import (
	"github.com/carloseduribeiro/auth-challenge/auth/configs"
	"github.com/carloseduribeiro/auth-challenge/auth/internal/infra/database"
	handlers2 "github.com/carloseduribeiro/auth-challenge/auth/internal/infra/http/handlers"
	"github.com/carloseduribeiro/auth-challenge/lib-utils/pkg/webserver"
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
	createUserHandler := handlers2.NewCreateUser(userRepository, uuid.NewUUID)
	if err = webServer.AddHandler(http.MethodPost, "/auth/users", createUserHandler.Handler); err != nil {
		log.Fatal(err)
	}

	sessionRepository := database.NewSessionRepository(dbPool)
	loginHandler := handlers2.NewLogin(userRepository, sessionRepository, config.JWTSecretKey, config.SessionMaxDuration)
	if err = webServer.AddHandler(http.MethodPost, "/auth/users/login", loginHandler.Handler); err != nil {
		log.Fatal(err)
	}

	webServer.Start()
}
