package model

import (
	"errors"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	Redis struct {
		Host     string `json:"host"`
		Password string `json:"password"`
		Port     string `json:"port"`
		Database int    `json:"db"`
	} `json:"redis"`
	Port              string `json:"port"`
	CheckJwt          bool   `json:"check_jwt"`
	SubAll            bool   `json:"sub_all"`
	UserAccess        string `json:"user_access"`
	BufferSize        int    `json:"buffer_size"`
	BroadcasterTimeout int    `json:"broadcaster_timeout_ms"`
	EnablePing        bool   `json:"enable_ping"`
}

type ConfigRequestApi struct {
	Version map[string]string `json:"version"`
	Domain  string            `json:"domain"`
	Login   struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"login"`
}

type ConfigRedisCache struct {
	Host     string `json:"host"`
	Password string `json:"password"`
	Port     string `json:"port"`
	Database int    `json:"db"`
}

func LoadConfig(file string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(file)
	v.SetConfigType("json")
	v.AutomaticEnv()

	// Set defaults
	v.SetDefault("buffer_size", 100)
	v.SetDefault("broadcaster_timeout_ms", 10)
	v.SetDefault("enable_ping", true)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return ValidateConfig(&config)
}

func ValidateConfig(config *Config) (*Config, error) {
	if config.Redis.Host == "" {
		return nil, errors.New("redis host is required")
	}
	if config.Redis.Port == "" {
		return nil, errors.New("redis port is required")
	}
	if _, err := strconv.Atoi(config.Port); err != nil {
		return nil, errors.New("invalid server port")
	}
	if config.BufferSize <= 0 {
		config.BufferSize = 100
	}
	if config.BroadcasterTimeout <= 0 {
		config.BroadcasterTimeout = 10
	}
	return config, nil
}