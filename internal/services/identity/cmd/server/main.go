package main

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/app"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`

	Database struct {
		Name     string `mapstructure:"name"`
		Database string `mapstructure:"database"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		SSLMode  string `mapstructure:"ssl_mode"`
	} `mapstructure:"database"`

	Security struct {
		Jwt struct {
			AccessSecret  string `mapstructure:"access_secret"`
			RefreshSecret string `mapstructure:"refresh_secret"`
			AccessTTL     string `mapstructure:"access_ttl"`
			RefreshTTL    string `mapstructure:"refresh_ttl"`
		} `mapstructure:"jwt"`
	} `mapstructure:"security"`
}

func main() {

	wd, err := os.Getwd()
	if err != nil {
		log.Printf("Failed to get working directory: %v", err)
		os.Exit(1)
	}

	// Load local environment variables
	loadEnv(path.Join(wd, ".env.local"))

	viper := viper.New()
	configPath := path.Join(wd, "internal/services/identity/config")

	viper.AddConfigPath(configPath)
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	// Bind Custom Environment Variables to Config keys
	viper.BindEnv("database.name", "IDENTITY_DB_NAME")
	viper.BindEnv("database.database", "IDENTITY_DB_NAME")
	viper.BindEnv("database.host", "IDENTITY_DB_HOST")
	viper.BindEnv("database.port", "IDENTITY_DB_PORT")
	viper.BindEnv("database.user", "IDENTITY_DB_USERNAME")
	viper.BindEnv("database.password", "IDENTITY_DB_PASSWORD")
	viper.BindEnv("database.ssl_mode", "IDENTITY_DB_SSL_MODE")

	viper.BindEnv("security.jwt.access_secret", "IDENTITY_JWT_ACCESS_SECRET")
	viper.BindEnv("security.jwt.refresh_secret", "IDENTITY_JWT_REFRESH_SECRET")
	viper.BindEnv("security.jwt.access_ttl", "IDENTITY_JWT_ACCESS_TTL")
	viper.BindEnv("security.jwt.refresh_ttl", "IDENTITY_JWT_REFRESH_TTL")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Failed to read config: %v", err)
		os.Exit(1)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Printf("Failed to unmarshal config: %v", err)
		os.Exit(1)
	}

	log.Print("test:", config.Server.Port)

	app, err := app.NewApplication()
	if err != nil {
		// log

		log.Printf("Failed to create application: %v", err)
		os.Exit(1)
	}

	if err := app.Run(); err != nil {
		// log

		log.Printf("Failed to run application: %v", err)
		os.Exit(1)
	}
}

func loadEnv(filepath string) {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return // Ignore if the file doesn't exist
	}
	lines := strings.Split(string(bytes), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			val = strings.Trim(val, `"'`)
			os.Setenv(key, val)
		}
	}
}
