package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
  SMTP SMTPConfig `mapstructure:"smtp"`
  Database DatabaseConfig `mapstructure:"database"`
  Server ServerConfig `mapstructure:"server"`
}

type SMTPConfig struct {
  Domain string `mapstructure:"domain"`
}

type DatabaseConfig struct {
  Host string `mapstructure:"host"`
  Port int `mapstructure:"port"`
  User string `mapstructure:"user"`
  Pass string `mapstructure:"pass"`
  Name string `mapstructure:"name"`
}

type ServerConfig struct {
  Port int `mapstructure:"port"`
}

func Load(path string) *Config {
  if path != "" {
    viper.SetConfigFile(path)
  } else {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AddConfigPath("/etc/postrun")
  }

  viper.SetEnvPrefix("POSTRUN")
  viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
  viper.AutomaticEnv()

  viper.SetDefault("server.port", 8080)

  if err := viper.ReadInConfig(); err != nil {
    if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
      log.Fatalf("config error: %v", err)
    }
    log.Println("no config file found, using defaults and env vars")
  }

  var cfg Config
  if err := viper.Unmarshal(&cfg); err != nil {
    log.Fatalf("config unmarshal error: %v", err)
  }

  return &cfg
}
