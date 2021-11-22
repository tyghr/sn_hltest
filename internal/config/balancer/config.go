package config

import (
	"github.com/spf13/viper"
)

const (
	TestMock   = "testmock"
	CacheRedis = "redis"
)

type Config struct {
	viper.Viper
	LogLevel int
	ApiPort  int

	ConsulServerAddr  string
	ConsulServiceName string
}

func NewConfig() *Config {
	return &Config{
		Viper:    *viper.New(),
		LogLevel: -1,
		ApiPort:  80,

		ConsulServerAddr:  "",
		ConsulServiceName: "",
	}
}

func (conf *Config) bindAllEnv() {
	_ = conf.BindEnv("loglevel", "SOCIAL_NETWORK_LOGLEVEL")
	_ = conf.BindEnv("apiport", "SOCIAL_NETWORK_APIPORT")

	_ = conf.BindEnv("consul_server_addr", "CONSUL_ADDR")
	_ = conf.BindEnv("consul_service_name", "CONSUL_SERVICE_NAME")
}

func (conf *Config) setDefaults() {
	conf.SetDefault("apiport", conf.ApiPort)
	conf.SetDefault("loglevel", conf.LogLevel)

	conf.SetDefault("consul_server_addr", conf.ConsulServerAddr)
	conf.SetDefault("consul_service_name", conf.ConsulServiceName)
}

//ReadSettings ...
// viper precedence order:
// 1 explicit call to Set
// 2 flag
// 3 env
// 4 config
// 5 key/value store
// 6 default
func (conf *Config) ReadAllSettings() error {
	conf.setDefaults()
	conf.bindAllEnv()

	conf.ApiPort = conf.GetInt("apiport")
	conf.LogLevel = conf.GetInt("loglevel")

	conf.ConsulServerAddr = conf.GetString("consul_server_addr")
	conf.ConsulServiceName = conf.GetString("consul_service_name")

	return nil
}
