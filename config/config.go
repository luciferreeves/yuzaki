package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DiscordToken   string
	DataSourceName string
	DatabaseDriver string
}

var validDBDrivers = map[string]bool{
	"sqlite":   true,
	"postgres": true,
	"mysql":    true,
	"mssql":    true,
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("loading env file: %w", err)
	}

	config := &Config{}
	requiredVars := map[string]*string{
		"DISCORD_BOT_TOKEN": &config.DiscordToken,
		"DSN":               &config.DataSourceName,
		"DATABASE_DRIVER":   &config.DatabaseDriver,
	}

	for envKey, configVar := range requiredVars {
		value, err := getEnv(envKey)
		if err != nil {
			return nil, fmt.Errorf("getting %s: %w", envKey, err)
		}
		*configVar = value
	}

	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	return config, nil
}

func validateConfig(config *Config) error {
	if config.DiscordToken == "" {
		return fmt.Errorf("discord token not found")
	}
	if config.DataSourceName == "" {
		return fmt.Errorf("data source name not found")
	}
	if config.DatabaseDriver == "" {
		return fmt.Errorf("database driver not found")
	}

	if !validDBDrivers[strings.ToLower(config.DatabaseDriver)] {
		return fmt.Errorf("invalid database driver: %s", config.DatabaseDriver)
	}

	return nil
}

func getEnv(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", fmt.Errorf("environment variable %s not found", key)
	}
	return strings.TrimSpace(value), nil
}

// Future implementation for getting other types of environment variables
// func getEnvInt(key string) (int, error) {
// 	value, err := getEnv(key)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return strconv.Atoi(value)
// }

// func getEnvBool(key string) (bool, error) {
// 	value, err := getEnv(key)
// 	if err != nil {
// 		return false, err
// 	}
// 	return strconv.ParseBool(value)
// }
