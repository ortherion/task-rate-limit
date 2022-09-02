package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

var cfg *Config

// New returns services config
func New() *Config {
	if cfg != nil {
		return cfg
	}

	return &Config{}
}

// Database - contains all parameters databases connection.
type Database struct {
	Host          string `yaml:"host"`
	Port          string `yaml:"port"`
	User          string `yaml:"user"`
	Password      string `yaml:"password"`
	UseMigrations string `yaml:"use_migrations"`
	Migrations    string `yaml:"migrations"`
	Name          string `yaml:"name"`
	SslMode       string `yaml:"sslmode"`
	Driver        string `yaml:"driver"`
	Timeout       int    `yaml:"timeout"`
}

// Rest - contains parameter rest json connection.
type Rest struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	DebugPort       int    `yaml:"debugPort"`
	ShutdownTimeout int    `yaml:"shutdownTimeout"`
	ReadTimeout     int    `yaml:"readTimeout"`
	WriteTimeout    int    `yaml:"writeTimeout"`
	IdleTimeout     int    `yaml:"idleTimeout"`
}

// App - contains all parameters project information.
type App struct {
	Debug       bool   `yaml:"debug"`
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
	Version     string `yaml:"version"`
}

// Metrics - contains all parameters metrics information.
type Metrics struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
	Path string `yaml:"path"`
}

// Jaeger - contains all parameters metrics information.
type Jaeger struct {
	Service string `yaml:"service"`
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
}

// Auth - contains parameter gRPC connection.
type Auth struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	MaxConnIdle int    `yaml:"maxConnectionIdle"`
	TimeOut     int    `yaml:"timeout"`
	MaxConnAge  int    `yaml:"maxConnectionAge"`
	SecretKey   string `yaml:"secretKey"`
}

type MsgHandlerConfig struct {
	RatePeriodMicroseconds int64 `yaml:"rate_period_seconds"`
	RequestsPerPeriod      int64 `yaml:"requests_per_period"`
}

// Config - contains all configuration parameters in config package.
type Config struct {
	App        App              `yaml:"app"`
	Rest       Rest             `yaml:"rest"`
	Database   Database         `yaml:"database"`
	Metrics    Metrics          `yaml:"metrics"`
	Jaeger     Jaeger           `yaml:"jaeger"`
	Auth       Auth             `yaml:"auth"`
	MsgHandler MsgHandlerConfig `yaml:"msgHandler"`
}

// ReadConfigYML - read configurations from file and init instance Config.
func ReadConfigYML(filePath string) error {
	if cfg != nil {
		return nil
	}

	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return err
	}

	return nil
}
