package configs

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

type Conf struct {
	DatabaseURL        string        `mapstructure:"DATABASE_URL"`
	DatabaseMaxConn    int32         `mapstructure:"DATABASE_MAX_CONN"`
	DatabaseMinConn    int32         `mapstructure:"DATABASE_MIN_CONN"`
	WebServerPort      string        `mapstructure:"WEB_SERVER_PORT"`
	JWTSecretKey       string        `mapstructure:"JWT_SECRET_KEY"`
	SessionMaxDuration time.Duration `mapstructure:"MAX_SESSION_MINUTES"`
}

func SetupDatabase(ctx context.Context, config *Conf) *pgxpool.Pool {
	dbConfig, err := pgxpool.ParseConfig(config.DatabaseURL)
	if err != nil {
		log.Fatalf("unable to parse database configuration: %v\n", err)
	}
	dbConfig.MaxConns = config.DatabaseMaxConn
	dbConfig.MinConns = config.DatabaseMinConn
	dbPool, err := pgxpool.New(ctx, config.DatabaseURL)
	if err != nil {
		log.Fatalf("unable to create connection pool: %v\n", err)
	}
	return dbPool
}
