package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const DefaultPath = "configs/local.yaml"

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Storage StorageConfig `yaml:"storage"`
}

type ServerConfig struct {
	Addr string `yaml:"addr"`
}

type StorageConfig struct {
	Mode        string `yaml:"mode"`
	DatabaseURL string `yaml:"database_url"`
}

func Default() Config {
	return Config{
		Server: ServerConfig{
			Addr: ":8080",
		},
		Storage: StorageConfig{
			Mode:        "memory",
			DatabaseURL: "postgres://postgres:postgres@localhost:5432/linkcuter?sslmode=disable",
		},
	}
}

func Load(path string) (Config, error) {
	if strings.TrimSpace(path) == "" {
		path = DefaultPath
	}

	cfg := Default()
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// если конфига нет — создаём дефолтный, чтобы было с чего стартовать
			if err := writeDefault(path, cfg); err != nil {
				return cfg, err
			}
			applyEnv(&cfg)
			return cfg, nil
		}
		return cfg, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	applyDefaults(&cfg)
	applyEnv(&cfg)
	return cfg, nil
}

func applyDefaults(cfg *Config) {
	def := Default()
	if strings.TrimSpace(cfg.Server.Addr) == "" {
		cfg.Server.Addr = def.Server.Addr
	}
	if strings.TrimSpace(cfg.Storage.Mode) == "" {
		cfg.Storage.Mode = def.Storage.Mode
	}
	if strings.TrimSpace(cfg.Storage.DatabaseURL) == "" {
		cfg.Storage.DatabaseURL = def.Storage.DatabaseURL
	}
}

func applyEnv(cfg *Config) {
	if v := os.Getenv("ADDR"); v != "" {
		cfg.Server.Addr = v
	}
	if v := os.Getenv("STORAGE"); v != "" {
		cfg.Storage.Mode = v
	}
	if v := os.Getenv("DATABASE_URL"); v != "" {
		cfg.Storage.DatabaseURL = v
	}
}

func writeDefault(path string, cfg Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}
