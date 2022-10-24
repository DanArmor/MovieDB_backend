package config

import "github.com/spf13/viper"

type Config struct {
	Port      string `mapstructure:"PORT"`
	SqlUrl    string `mapstructure:"SqlUrl"`
	JWTsecret string `mapstructure:"JWT_SECRET_KEY"`
}

func LoadConfig() (c Config, err error) {
	viper.AddConfigPath("./pkg/config/envs")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&c)

	return
}
