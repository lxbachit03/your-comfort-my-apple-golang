package config

import (
	"github.com/spf13/viper"
)

type ServerConfig struct {
	AppEnv string `mapstructure:"app_env"`
	Port   int    `mapstructure:"port"`
}

type DatabaseConfig struct {
	Name     string `mapstructure:"name"`
	Database string `mapstructure:"database"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

type JWTConfig struct {
	Secret          string `mapstructure:"secret"`
	RefreshSecret   string `mapstructure:"refresh_secret"`
	EncryptionKey   string `mapstructure:"encryption_key"`
	AccessTokenTtl  int    `mapstructure:"access_ttl"`
	RefreshTokenTtl int    `mapstructure:"refresh_ttl"`
	Issuer          string `mapstructure:"issuer"`
}

type SecurityConfig struct {
	Jwt JWTConfig `mapstructure:"jwt"`
}

type CacheConfig struct {
	Redis RedisConfig `mapstructure:"redis"`
}

type RedisConfig struct {
	DB       int    `mapstructure:"db"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Security SecurityConfig `mapstructure:"security"`
	Cache    CacheConfig    `mapstructure:"cache"`
}

var AppConfig Config

func LoadConfig(configPath string, env string) error {
	viper := viper.New()

	viper.AddConfigPath(configPath)
	viper.SetConfigName(env)
	viper.SetConfigType("yaml")

	// Default values if env variables or config file not provided
	// AppConfig = Config{
	// 	Server: ServerConfig{
	// 		AppEnv: "local",
	// 		Port:   8080,
	// 	},
	// 	Database: DatabaseConfig{
	// 		Name:     "identity",
	// 		Database: "identity",
	// 		Host:     "localhost",
	// 		Port:     5432,
	// 		User:     "admin",
	// 		Password: "adminpassword",
	// 		SSLMode:  "disable",
	// 	},
	// 	Security: SecurityConfig{
	// 		Jwt: JWTConfig{
	// 			Secret:          "your_jwt_access_secret_key",
	// 			RefreshSecret:   "your_jwt_refresh_secret_key",
	// 			AccessTokenTtl:  60 * 60 * 24 * 7,
	// 			RefreshTokenTtl: 60 * 60 * 24 * 30,
	// 			Issuer:          "your_jwt_issuer",
	// 		},
	// 	},
	// 	Cache: CacheConfig{
	// 		Redis: RedisConfig{
	// 			DB:       0,
	// 			Host:     "127.0.0.1",
	// 			Port:     6379,
	// 			User:     "",
	// 			Password: "",
	// 		},
	// 	},
	// }

	// Bind Custom Environment Variables to Config keys
	viper.BindEnv("server.app_env", "APP_ENV")
	viper.BindEnv("server.port", "PORT")

	viper.BindEnv("database.name", "IDENTITY_POSTGRES_DB_NAME")
	viper.BindEnv("database.database", "IDENTITY_POSTGRES_DB_NAME")
	viper.BindEnv("database.host", "IDENTITY_POSTGRES_DB_HOST")
	viper.BindEnv("database.port", "IDENTITY_POSTGRES_DB_PORT")
	viper.BindEnv("database.user", "IDENTITY_POSTGRES_DB_USER")
	viper.BindEnv("database.password", "IDENTITY_POSTGRES_DB_PASSWORD")
	viper.BindEnv("database.ssl_mode", "IDENTITY_POSTGRES_DB_SSL_MODE")

	viper.BindEnv("security.jwt.secret", "IDENTITY_JWT_SECRET")
	viper.BindEnv("security.jwt.refresh_secret", "IDENTITY_JWT_REFRESH_SECRET")
	viper.BindEnv("security.jwt.encryption_key", "IDENTITY_JWT_ENCRYPTION_KEY")
	viper.BindEnv("security.jwt.access_ttl", "IDENTITY_JWT_ACCESS_TOKEN_TTL")
	viper.BindEnv("security.jwt.refresh_ttl", "IDENTITY_JWT_REFRESH_TOKEN_TTL")
	viper.BindEnv("security.jwt.issuer", "IDENTITY_JWT_ISSUER")

	viper.BindEnv("cache.redis.db", "CACHE_REDIS_DB")
	viper.BindEnv("cache.redis.host", "CACHE_REDIS_HOST")
	viper.BindEnv("cache.redis.port", "CACHE_REDIS_PORT")
	viper.BindEnv("cache.redis.user", "CACHE_REDIS_USER")
	viper.BindEnv("cache.redis.password", "CACHE_REDIS_PASSWORD")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		return err
	}

	return nil
}
