package config

import (
	"os"
	"time"

	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Port int `yaml:"port"`
}

type EtcdConfig struct {
	Endpoints         []string `yaml:"endpoints"`
	DialTimeoutSecond int      `yaml:"dial_timeout_seconds"`
}

type LogConfig struct {
	Level string `yaml:"level"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
	Etcd   EtcdConfig   `yaml:"etcd"`
	Log    LogConfig    `yaml:"log"`
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Validate/Set defaults
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if len(cfg.Etcd.Endpoints) == 0 {
		cfg.Etcd.Endpoints = []string{"localhost:2379"}
	}
	if cfg.Etcd.DialTimeoutSecond == 0 {
		cfg.Etcd.DialTimeoutSecond = 5
	}
	if cfg.Log.Level == "" {
		cfg.Log.Level = "info"
	}

	return &cfg, nil
}

func (c *Config) EtcdDialTimeout() time.Duration {
	return time.Duration(c.Etcd.DialTimeoutSecond) * time.Second
}
