package config

import (
	"bytes"
	"strings"

	_ "embed"

	"github.com/spf13/viper"
)

const (
	EnvDev    = "dev"
	EnvUat    = "uat"
	EnvLambda = "lambda"
	EnvLocal  = "local"
)

//go:embed config.yaml
var defaultConfig []byte

var NonProductionEnvironments = []string{EnvDev, EnvUat, EnvLambda, EnvLocal}

type (
	Config struct {
		Env             string            `mapstructure:"env"`
		App             string            `mapstructure:"app"`
		LogLevel        string            `mapstructure:"log_level"`
		DBLog           bool              `mapstructure:"db_log"`
		HttpAddress     uint32            `mapstructure:"http_address"`
		AllowOrigins    string            `mapstructure:"allow_origins"`
		SensitiveFields map[string]string `mapstructure:"sensitive_fields"`
		Database        databaseConfig    `mapstructure:"database"`
		AutoMigration   bool              `mapstructure:"auto_migration"`
		Cache           cacheConfig       `mapstructure:"redis_cache"`
		JwtConfig       jwtConfig         `mapstructure:"jwt_config"`
		HmacConfig      hmacConfig        `mapstructure:"hmac_config"`
		Timezone        string            `mapstructure:"timezone"`
	}

	hmacConfig struct {
		SecretKey     string `mapstructure:"secret_key"`
		MaxAgeSeconds int    `mapstructure:"max_age_seconds"`
	}

	databaseConfig struct {
		Host       string `mapstructure:"host"`
		Port       int    `mapstructure:"port"`
		Username   string `mapstructure:"username"`
		Password   string `mapstructure:"password"`
		Database   string `mapstructure:"database"`
		SearchPath string `mapstructure:"search_path"`
		Timezone   string `mapstructure:"timezone"`
		SSLMode    string `json:"ssl_mode" mapstructure:"ssl_mode"`
	}

	cacheConfig struct {
		Host     string `mapstructure:"host"`
		Port     uint32 `mapstructure:"port"`
		Password string `mapstructure:"password"`
	}

	jwtConfig struct {
		SecretKey       string `mapstructure:"secret_key"`
		Duration        int    `mapstructure:"duration"`         // in seconds
		DurationRefresh int    `mapstructure:"duration_refresh"` // in seconds
		Algorithm       string `mapstructure:"algorithm"`
		SessionDuration int    `mapstructure:"session_duration"`
	}

	adapterConfig struct {
		BaseUrl                     string `mapstructure:"base_url"`
		Username                    string `mapstructure:"username"`
		Password                    string `mapstructure:"password"`
		AuthUrl                     string `mapstructure:"auth_url"`
		ApiKey                      string `mapstructure:"api_key"`
		ClientID                    string `mapstructure:"client_id"`
		ClientSecret                string `mapstructure:"client_secret"`
		SecretKey                   string `mapstructure:"secret_key"`
		Timeout                     int    `mapstructure:"timeout"`
		EmitterBufferSize           int    `mapstructure:"emitter_buffer_size"`
		EmitterFlushIntervalSeconds int    `mapstructure:"emitter_flush_interval_seconds"`
		AccessKey                   string `mapstructure:"access_key"`
	}
)

func Load() (*Config, error) {
	var cfg = &Config{}

	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))

	viper.AutomaticEnv()

	err := viper.ReadConfig(bytes.NewBuffer(defaultConfig))
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil

}
