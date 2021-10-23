package chat

import (
	"github.com/spf13/viper"
)

const (
	TestMock   = "testmock"
	CacheRedis = "redis"
)

type Config struct {
	viper.Viper
	// ConfigName string
	// ConfigType string
	// ConfigPath string
	LogLevel           int
	ApiPort            int
	SessionValidateUrl string

	CacheType      string
	CacheNodes     []string
	CachePass      string
	CacheClustered bool
}

func NewConfig() *Config {
	return &Config{
		Viper: *viper.New(),
		// ConfigName: "config",
		// ConfigType: "json",
		// ConfigPath: "./", // change in prod to "/etc/social_network/"
		LogLevel:           -1,
		ApiPort:            80,
		SessionValidateUrl: "http://127.0.0.1/session_validate",

		CacheType:      CacheRedis,
		CacheNodes:     []string{"redis_node_0:6379", "redis_node_1:6379", "redis_node_2:6379", "redis_node_3:6379", "redis_node_4:6379", "redis_node_5:6379"},
		CachePass:      "testpass",
		CacheClustered: false,
	}
}

func (conf *Config) bindAllEnv() {
	_ = conf.BindEnv("loglevel", "SOCIAL_NETWORK_LOGLEVEL")
	_ = conf.BindEnv("apiport", "SOCIAL_NETWORK_APIPORT")
	_ = conf.BindEnv("session_validate_url", "SOCIAL_NETWORK_SESSION_VALIDATE_URL")

	_ = conf.BindEnv("cachetype", "SOCIAL_NETWORK_CACHETYPE")
	_ = conf.BindEnv("cachenodes", "SOCIAL_NETWORK_CACHENODES")
	_ = conf.BindEnv("cachepass", "SOCIAL_NETWORK_CACHEPASS")
	_ = conf.BindEnv("cacheclustered", "SOCIAL_NETWORK_CACHECLUSTERED")
}

func (conf *Config) setDefaults() {
	conf.SetDefault("apiport", conf.ApiPort)
	conf.SetDefault("loglevel", conf.LogLevel)
	conf.SetDefault("session_validate_url", conf.SessionValidateUrl)

	conf.SetDefault("cachetype", conf.CacheType)
	conf.SetDefault("cachenodes", conf.CacheNodes)
	conf.SetDefault("cachepass", conf.CachePass)
	conf.SetDefault("cacheclustered", conf.CacheClustered)
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
	conf.SessionValidateUrl = conf.GetString("session_validate_url")

	conf.CacheType = conf.GetString("cachetype")
	conf.CacheNodes = conf.GetStringSlice("cachenodes")
	conf.CachePass = conf.GetString("cachepass")
	conf.CacheClustered = conf.GetBool("cacheclustered")

	return nil
}
