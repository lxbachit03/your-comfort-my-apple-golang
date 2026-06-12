package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/config"
)

// var DB sqlc.Querier
var DBPool *pgxpool.Pool

func InitDB() error {
	var connStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.AppConfig.Database.Host,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.User,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.Database,
		config.AppConfig.Database.SSLMode)

	conf, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("error parsing DB config: %v", err)
	}

	conf.MaxConns = 50
	conf.MinConns = 5
	conf.MaxConnLifetime = 30 * time.Minute
	conf.MaxConnIdleTime = 5 * time.Minute
	conf.HealthCheckPeriod = 1 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	DBPool, err = pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return fmt.Errorf("error creating DB pool: %v", err)
	}

	// DB = sqlc.New(DBPool)

	if err := DBPool.Ping(ctx); err != nil {
		return fmt.Errorf("db ping error: %v", err)
	}

	log.Printf("✅ Connected Database Postgresql")

	return nil
}
