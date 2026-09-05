package config

import(
	"fmt"
	"time"
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT JWTConfig `mapstructure:"jwt"`
}

type ServerConfig struct {
    Port int `mapstructure:"port"`
	ReadTimeout time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	Host string `mapstructure:"host"`
	Port int `mapstructure:"port"`
	User string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name string `mapstructure:"name"`
	SSLMode string `mapstructure:"sslmode"`
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", /*what is this?*/
		d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode)
}

type JWTConfig struct {
	Secret         string        `mapstructure:"secret"`
	AccessTokenTTL time.Duration `mapstructure:"access_token_ttl"`
}

func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.AutomaticEnv() // env vars override yaml, e.g. JWT_SECRET

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshalling config: %w", err)
	}
	return &cfg, nil
}