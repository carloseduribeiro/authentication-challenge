package configs

import "github.com/spf13/viper"

type Conf struct {
	DatabaseURL     string `mapstructure:"DATABASE_URL"`
	DatabaseMaxConn int32  `mapstructure:"DATABASE_MAX_CONN"`
	DatabaseMinConn int32  `mapstructure:"DATABASE_MIN_CONN"`
	WebServerPort   string `mapstructure:"WEB_SERVER_PORT"`
}

func LoadConfig(path string) (*Conf, error) {
	var cfg *Conf
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
