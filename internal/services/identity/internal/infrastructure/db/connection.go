package db

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/logger"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/pgx"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/config"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/infrastructure/db/sqlc"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/infrastructure/utils"
)

var DB sqlc.Querier
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

	// Init SQL logger
	rootDir := utils.GetWorkingDir()
	sqlLoggerPath := path.Join(rootDir, "logs/identity/sql.log")
	sqlLogger := logger.NewLogger(logger.LoggerConfig{
		Level:      "info",
		Filename:   sqlLoggerPath,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     5,
		Compress:   true,
		AppEnv:     config.AppConfig.Server.AppEnv,
	})

	conf.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger: &pgx.PgxLogTracer{
			Logger:         *sqlLogger,
			SlowQueryLimit: 500 * time.Millisecond,
		},
		LogLevel: tracelog.LogLevelDebug,
	}

	// connection pool settings
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

	// Init sqlc
	DB = sqlc.New(DBPool)

	if err := DBPool.Ping(ctx); err != nil {
		return fmt.Errorf("db ping error: %v", err)
	}

	logger.Log.Info().Msg("✅ Connected Database Postgresql")

	return nil
}
