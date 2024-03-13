package config

import "github.com/spf13/viper"

func LoadConfig[T any](path string) (*T, error) {
	var cfg *T
	viper.AutomaticEnv()
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
