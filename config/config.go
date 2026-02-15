package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Logger   LoggerConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int32
	MinConns int32
}

type JWTConfig struct {
	Secret         string
	ExpirationHour int
}

type LoggerConfig struct {
	Level string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "password")
	viper.SetDefault("DB_NAME", "bidding")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("DB_MAX_CONNS", 25)
	viper.SetDefault("DB_MIN_CONNS", 5)
	viper.SetDefault("JWT_SECRET", "your-secret-key-change-in-production")
	viper.SetDefault("JWT_EXPIRATION_HOUR", 24)
	viper.SetDefault("LOG_LEVEL", "info")

	cfg := &Config{
		Server: ServerConfig{
			Port: viper.GetString("SERVER_PORT"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
			MaxConns: viper.GetInt32("DB_MAX_CONNS"),
			MinConns: viper.GetInt32("DB_MIN_CONNS"),
		},
		JWT: JWTConfig{
			Secret:         viper.GetString("JWT_SECRET"),
			ExpirationHour: viper.GetInt("JWT_EXPIRATION_HOUR"),
		},
		Logger: LoggerConfig{
			Level: viper.GetString("LOG_LEVEL"),
		},
	}

	log.Printf("Configuration loaded successfully")
	return cfg, nil
}

// GetDSN returns PostgreSQL connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}