package configs

import "github.com/spf13/viper"

type conf struct {
	dBDriver          string `mapstructure:"DB_DRIVER"`
	dBHost            string `mapstructure:"DB_HOST"`
	dBPort            string `mapstructure:"DB_PORT"`
	dBUser            string `mapstructure:"DB_USER"`
	dBPassword        string `mapstructure:"DB_PASSWORD"`
	dBName            string `mapstructure:"DB_NAME"`
	webServerPort     string `mapstructure:"WEB_SERVER_PORT"`
	gRPCServerPort    string `mapstructure:"GRPC_SERVER_PORT"`
	graphQLServerPort string `mapstructure:"GRAPHQL_SERVER_PORT"`
}

func (c *conf) GetDBDriver() string {
	return c.dBDriver
}

func (c *conf) GetDBHost() string {
	return c.dBHost
}

func (c *conf) GetDBPort() string {
	return c.dBPort
}

func (c *conf) GetDBUser() string {
	return c.dBUser
}

func (c *conf) GetDBPassword() string {
	return c.dBPassword
}

func (c *conf) GetDBName() string {
	return c.dBName
}

func (c *conf) GetWebServerPort() string {
	return c.webServerPort
}

func (c *conf) GetGRPCServerPort() string {
	return c.gRPCServerPort
}

func (c *conf) GetGraphQLServerPort() string {
	return c.graphQLServerPort
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	cfg.dBDriver = viper.GetString("DB_DRIVER")
	cfg.dBHost = viper.GetString("DB_HOST")
	cfg.dBPort = viper.GetString("DB_PORT")
	cfg.dBUser = viper.GetString("DB_USER")
	cfg.dBPassword = viper.GetString("DB_PASSWORD")
	cfg.dBName = viper.GetString("DB_NAME")
	cfg.webServerPort = viper.GetString("WEB_SERVER_PORT")
	cfg.gRPCServerPort = viper.GetString("GRPC_SERVER_PORT")
	cfg.graphQLServerPort = viper.GetString("GRAPHQL_SERVER_PORT")

	return cfg, err
}
