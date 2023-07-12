package config

import (
	"os"

	"github.com/spf13/viper"
)

const (
	DevEnvironment = "DEV"
)

type AppConfig struct {
	App struct {
		Environment string
	}

	ByBit struct {
		APIKey    string
		APISecret string
		BaseURL   string
	}
}

var (
	cfg *AppConfig
)

func Config() *AppConfig {
	if cfg == nil {
		loadConfig()
	}

	return cfg
}

func loadConfig() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Ignore config file not found, perhaps we will use environment variables.
	_ = viper.ReadInConfig()

	cfg = &AppConfig{}

	// App.
	cfg.App.Environment = os.Getenv("ENVIRONMENT")

	// ByBit.
	cfg.ByBit.APIKey = viper.GetString("BYBIT_API_KEY")
	cfg.ByBit.APISecret = viper.GetString("BYBIT_API_SECRET")
	cfg.ByBit.BaseURL = viper.GetString("BYBIT_BASE_URL")
}
