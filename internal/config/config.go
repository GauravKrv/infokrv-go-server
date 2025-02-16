// internal/config/config.go
package config

import (
    "os"
    "time"

    "github.com/spf13/viper"
)
type Config struct {
    Server struct {
        Port         string        `mapstructure:"port"`
        ReadTimeout  time.Duration `mapstructure:"read_timeout"`
        WriteTimeout time.Duration `mapstructure:"write_timeout"`
    } `mapstructure:"server"`

    MongoDB struct {
        URI      string `mapstructure:"uri"`
        Database string `mapstructure:"database"`
    } `mapstructure:"mongodb"`

    JWT struct {
        Secret string        `mapstructure:"secret"`
        TTL    time.Duration `mapstructure:"ttl"`
    } `mapstructure:"jwt"`
}

func Load() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AddConfigPath("./configs")

    viper.AutomaticEnv()

    var config Config
    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }
    // Override MongoDB URI from environment variable if present
    if mongoURI := os.Getenv("MONGODB_URI"); mongoURI != "" {
            viper.Set("mongodb.uri", mongoURI)
    }

    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }
    return &config, nil
}