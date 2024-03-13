package configs

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
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

func LoadConfig(path string) (*Conf, error) {
	var cfg *Conf
	viper.AutomaticEnv()
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
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

func RunMigrations(databaseURL string) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatalf("error open connection to apply migration: %s", err)
	}
	defer db.Close()
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("could not init driver: %s", err)
	}
	defer driver.Close()
	m, err := migrate.NewWithDatabaseInstance("file://../internal/infra/database/migrations", "pgx", driver)
	if err != nil {
		log.Fatalf("could not apply the migration: %s", err)
	}
	m.Up()
	defer m.Close()
}
