package config

import (
	"errors"

	"lim/internal/lim-repo/flags"

	"github.com/spf13/viper"
)

type Config struct {
	Postgres      flags.PostgresFlags `mapstructure:"postgres"`
	Address       string              `mapstructure:"address"`
	Port          string              `mapstructure:"port"`
	LogLevel      string              `mapstructure:"loglevel"`
	LogFile       string              `mapstructure:"logfile"`
	Mode          string              `mapstructure:"mode"`
	AdminLogin    string              `mapstructure:"admin_login"`
	AdminPassword string              `mapstructure:"admin_password"`

	FirstTrainingTime int `mapstructure:"first_training_time"`
	LastTrainingTime  int `mapstructure:"last_training_time"`
}

func (c *Config) ParseConfig(configFileName, pathToConfig string) error {
	v := viper.New()
	v.SetConfigName(configFileName)
	v.SetConfigType("json")
	v.AddConfigPath(pathToConfig)

	err := v.ReadInConfig()
	if err != nil {
		return err
	}

	err = v.Unmarshal(c)
	if err != nil {
		return err
	}

	var emptyPostgres flags.PostgresFlags
	if c.Postgres == emptyPostgres {
		return errors.New("no postgres info in config")
	}

	if c.Postgres.Host == "" {
		return errors.New("no postgres host in config")
	}
	if c.Postgres.DBName == "" {
		return errors.New("no postgres db name in config")
	}
	if c.Postgres.User == "" {
		return errors.New("no postgres db user in config")
	}
	if c.Postgres.Password == "" {
		return errors.New("no postgres db password in config")
	}
	if c.Postgres.Port == "" {
		return errors.New("no postgres db port in config")
	}

	return nil
}
