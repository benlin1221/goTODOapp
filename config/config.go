package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper
	database *DatabaseConfig
}

var instantiated *Config
var once sync.Once

func GetInstance() *Config {
	once.Do(func() {
		instantiated = &Config{
			Viper: viper.New(),
		}

		// Set default configurations
		instantiated.setDefaults()

		// Select the .env file
		instantiated.SetConfigName(".env")
		instantiated.SetConfigType("dotenv")
		instantiated.AddConfigPath(".")

		// Automatically refresh environment variables
		instantiated.AutomaticEnv()

		// Read configuration
		if err := instantiated.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				fmt.Println("failed to read configuration:", err.Error())
				os.Exit(1)
			}
		}

		instantiated.setDatabaseConfig()
	})
	return instantiated
}

func (config *Config) setDefaults() {
	// Set default database configuration
	config.SetDefault("DB_DRIVER", "postgresql")
	config.SetDefault("DB_HOST", "localhost")
	config.SetDefault("DB_USERNAME", "")
	config.SetDefault("DB_PASSWORD", "")
	config.SetDefault("DB_PORT", 5432)
	config.SetDefault("DB_NAME", "maintainer")
}

func (config *Config) setDatabaseConfig() {
	config.database = &DatabaseConfig{
		Default: DatabaseDriver{
			Driver:   config.GetString("DB_DRIVER"),
			Host:     config.GetString("DB_HOST"),
			Username: config.GetString("DB_USERNAME"),
			Password: config.GetString("DB_PASSWORD"),
			DBName:   config.GetString("DB_NAME"),
			Port:     config.GetInt("DB_PORT"),
		},
	}
}

func (config *Config) getDatabaseConfig() *DatabaseConfig {
	return config.database
}
