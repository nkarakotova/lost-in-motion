package config

import (
	"lim/internal/lim-repo/flags"
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
