package config

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port         int           `yaml:"port"`
	DataDir      string        `yaml:"data_dir"`
	DBPath       string        `yaml:"db_path"`
	ServeUI      bool          `yaml:"serve_ui"`
	UITargetDir  string        `yaml:"ui_dist"`
	SessionTTL   time.Duration `yaml:"session_ttl"`
	TokenTTL     time.Duration `yaml:"token_ttl"`
	CSRFEnabled  bool          `yaml:"csrf_enabled"`
	SecretKey    string        `yaml:"secret_key"`
	TLSCertPath  string        `yaml:"tls_cert"`
	TLSKeyPath   string        `yaml:"tls_key"`
	CACertPath   string        `yaml:"ca_cert"`
	AdminEmail   string        `yaml:"admin_email"`
	AdminPass    string        `yaml:"admin_password"`
	PoolGRPCAddr string        `yaml:"pool_grpc_addr"`
}

type fileConfig struct {
	Port         int    `yaml:"port"`
	DataDir      string `yaml:"data_dir"`
	DBPath       string `yaml:"db_path"`
	ServeUI      *bool  `yaml:"serve_ui"`
	UITargetDir  string `yaml:"ui_dist"`
	SessionTTL   string `yaml:"session_ttl"`
	TokenTTL     string `yaml:"token_ttl"`
	CSRFEnabled  *bool  `yaml:"csrf_enabled"`
	SecretKey    string `yaml:"secret_key"`
	TLSCertPath  string `yaml:"tls_cert"`
	TLSKeyPath   string `yaml:"tls_key"`
	CACertPath   string `yaml:"ca_cert"`
	AdminEmail   string `yaml:"admin_email"`
	AdminPass    string `yaml:"admin_password"`
	PoolGRPCAddr string `yaml:"pool_grpc_addr"`
}

func Load() (*Config, error) {
	cfg := &Config{
		DataDir:      filepath.Clean("../backend/data"),
		DBPath:       filepath.Clean("../backend/data/openaction.db"),
		ServeUI:      true,
		UITargetDir:  filepath.Clean("../backend/web/dist"),
		SessionTTL:   7 * 24 * time.Hour,
		TokenTTL:     90 * 24 * time.Hour,
		CSRFEnabled:  true,
		AdminEmail:   "admin@openaction.local",
		AdminPass:    "admin123",
		PoolGRPCAddr: ":7443",
	}

	if filePath := os.Getenv("OA_CONFIG"); filePath != "" {
		if err := applyFileConfig(cfg, filePath); err != nil {
			return nil, err
		}
	} else if _, err := os.Stat("config.yaml"); err == nil {
		if err := applyFileConfig(cfg, "config.yaml"); err != nil {
			return nil, err
		}
	}

	if port := os.Getenv("OA_PORT"); port != "" {
		value, err := strconv.Atoi(port)
		if err != nil {
			return nil, errors.New("OA_PORT must be integer")
		}
		cfg.Port = value
	}
	if cfg.Port == 0 {
		return nil, errors.New("OA_PORT is required")
	}

	if dir := os.Getenv("OA_DATA_DIR"); dir != "" {
		cfg.DataDir = filepath.Clean(dir)
	}
	if db := os.Getenv("OA_DB_PATH"); db != "" {
		cfg.DBPath = filepath.Clean(db)
	}
	if ui := os.Getenv("OA_UI_DIST"); ui != "" {
		cfg.UITargetDir = filepath.Clean(ui)
	}
	if serve := os.Getenv("OA_SERVE_UI"); serve != "" {
		cfg.ServeUI = serve == "1" || serve == "true"
	}
	if ttl := os.Getenv("OA_SESSION_TTL"); ttl != "" {
		if parsed, err := time.ParseDuration(ttl); err == nil {
			cfg.SessionTTL = parsed
		}
	}
	if ttl := os.Getenv("OA_TOKEN_TTL"); ttl != "" {
		if parsed, err := time.ParseDuration(ttl); err == nil {
			cfg.TokenTTL = parsed
		}
	}
	if csrf := os.Getenv("OA_CSRF"); csrf != "" {
		cfg.CSRFEnabled = csrf == "1" || csrf == "true"
	}
	if v := os.Getenv("OA_SECRET_KEY"); v != "" {
		cfg.SecretKey = v
	}
	if v := os.Getenv("OA_TLS_CERT"); v != "" {
		cfg.TLSCertPath = v
	}
	if v := os.Getenv("OA_TLS_KEY"); v != "" {
		cfg.TLSKeyPath = v
	}
	if v := os.Getenv("OA_CA_CERT"); v != "" {
		cfg.CACertPath = v
	}
	if v := os.Getenv("OA_ADMIN_EMAIL"); v != "" {
		cfg.AdminEmail = v
	}
	if v := os.Getenv("OA_ADMIN_PASSWORD"); v != "" {
		cfg.AdminPass = v
	}
	if v := os.Getenv("OA_POOL_GRPC_ADDR"); v != "" {
		cfg.PoolGRPCAddr = v
	}
	if cfg.SecretKey == "" {
		return nil, errors.New("OA_SECRET_KEY is required")
	}

	return cfg, nil
}

func applyFileConfig(cfg *Config, path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var fc fileConfig
	if err := yaml.Unmarshal(raw, &fc); err != nil {
		return err
	}
	if fc.Port != 0 {
		cfg.Port = fc.Port
	}
	if fc.DataDir != "" {
		cfg.DataDir = filepath.Clean(fc.DataDir)
	}
	if fc.DBPath != "" {
		cfg.DBPath = filepath.Clean(fc.DBPath)
	}
	if fc.ServeUI != nil {
		cfg.ServeUI = *fc.ServeUI
	}
	if fc.UITargetDir != "" {
		cfg.UITargetDir = filepath.Clean(fc.UITargetDir)
	}
	if fc.SessionTTL != "" {
		if parsed, err := time.ParseDuration(fc.SessionTTL); err == nil {
			cfg.SessionTTL = parsed
		}
	}
	if fc.TokenTTL != "" {
		if parsed, err := time.ParseDuration(fc.TokenTTL); err == nil {
			cfg.TokenTTL = parsed
		}
	}
	if fc.CSRFEnabled != nil {
		cfg.CSRFEnabled = *fc.CSRFEnabled
	}
	if fc.SecretKey != "" {
		cfg.SecretKey = fc.SecretKey
	}
	if fc.TLSCertPath != "" {
		cfg.TLSCertPath = fc.TLSCertPath
	}
	if fc.TLSKeyPath != "" {
		cfg.TLSKeyPath = fc.TLSKeyPath
	}
	if fc.CACertPath != "" {
		cfg.CACertPath = fc.CACertPath
	}
	if fc.AdminEmail != "" {
		cfg.AdminEmail = fc.AdminEmail
	}
	if fc.AdminPass != "" {
		cfg.AdminPass = fc.AdminPass
	}
	if fc.PoolGRPCAddr != "" {
		cfg.PoolGRPCAddr = fc.PoolGRPCAddr
	}

	return nil
}
