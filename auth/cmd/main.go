package main

import (
	"github.com/carloseduribeiro/auth-challenge/auth/configs"
	"github.com/carloseduribeiro/auth-challenge/auth/internal/infra/database"
	"github.com/carloseduribeiro/auth-challenge/auth/internal/infra/http/handlers"
	"github.com/carloseduribeiro/auth-challenge/lib-utils/pkg/config"
	"github.com/carloseduribeiro/auth-challenge/lib-utils/pkg/db"
	"github.com/carloseduribeiro/auth-challenge/lib-utils/pkg/webserver"
	"github.com/google/uuid"
	"golang.org/x/net/context"
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

	ctx := context.TODO()
	dbPool := configs.SetupDatabase(ctx, cfg)
	defer dbPool.Close()

	webServer := webserver.NewWebServer(cfg.WebServerPort)

	userRepository := database.NewUserRepository(dbPool)
	createUserHandler := handlers.NewCreateUser(userRepository, uuid.NewUUID)
	if err = webServer.AddHandler(http.MethodPost, "/auth/users", createUserHandler.Handler); err != nil {
		log.Fatal(err)
	}

	sessionRepository := database.NewSessionRepository(dbPool)
	loginHandler := handlers.NewLogin(userRepository, sessionRepository, cfg.JWTSecretKey, cfg.SessionMaxDuration)
	if err = webServer.AddHandler(http.MethodPost, "/auth/users/login", loginHandler.Handler); err != nil {
		log.Fatal(err)
	}

	webServer.Start()
}
