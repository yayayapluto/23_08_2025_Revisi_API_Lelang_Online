package utils

import (
	"github.com/spf13/viper"
)

type Env struct {
	DBHOST              string
	DBUSER              string
	DBPASSWORD          string
	DBNAME              string
	DBPORT              string
	JWT_SECRET          string `mapstructure:"JWT_SECRET"`
	MIDTRANS_SERVER_KEY string `mapstructure:"MIDTRANS_SERVER_KEY"`
	MIDTRANS_MODE       string `mapstructure:"MIDTRANS_MODE"`
}

func LoadEnv() (*Env, error) {
	var env Env

	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&env); err != nil {
		return nil, err
	}

	return &env, nil
}
