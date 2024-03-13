package configs

type Conf struct {
	DatabaseURL     string `mapstructure:"DATABASE_URL"`
	DatabaseMaxConn int32  `mapstructure:"DATABASE_MAX_CONN"`
	DatabaseMinConn int32  `mapstructure:"DATABASE_MIN_CONN"`
	WebServerPort   string `mapstructure:"WEB_SERVER_PORT"`
}
