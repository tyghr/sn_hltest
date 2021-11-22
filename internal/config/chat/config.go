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

	ConsulServerAddr  string
	ConsulServiceName string
	ConsulServiceID   string
	ConsulAgentAddr   string
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

		ConsulServerAddr:  "",
		ConsulServiceName: "",
		ConsulServiceID:   "",
		ConsulAgentAddr:   "",
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

	_ = conf.BindEnv("consul_server_addr", "CONSUL_ADDR")
	_ = conf.BindEnv("consul_service_name", "CONSUL_SERVICE_NAME")
	_ = conf.BindEnv("consul_service_id", "CONSUL_SERVICE_ID")
	_ = conf.BindEnv("consul_agent_addr", "CONSUL_AGENT_ADDR")
}

func (conf *Config) setDefaults() {
	conf.SetDefault("apiport", conf.ApiPort)
	conf.SetDefault("loglevel", conf.LogLevel)
	conf.SetDefault("session_validate_url", conf.SessionValidateUrl)

	conf.SetDefault("cachetype", conf.CacheType)
	conf.SetDefault("cachenodes", conf.CacheNodes)
	conf.SetDefault("cachepass", conf.CachePass)
	conf.SetDefault("cacheclustered", conf.CacheClustered)

	conf.SetDefault("consul_server_addr", conf.ConsulServerAddr)
	conf.SetDefault("consul_service_name", conf.ConsulServiceName)
	conf.SetDefault("consul_service_id", conf.ConsulServiceID)
	conf.SetDefault("consul_agent_addr", conf.ConsulAgentAddr)
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

	conf.ConsulServerAddr = conf.GetString("consul_server_addr")
	conf.ConsulServiceName = conf.GetString("consul_service_name")
	conf.ConsulServiceID = conf.GetString("consul_service_id")
	conf.ConsulAgentAddr = conf.GetString("consul_agent_addr")

	return nil
}
