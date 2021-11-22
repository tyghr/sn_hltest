package config

import (
	"os"

	"github.com/spf13/viper"
)

const (
	TestMock        = "testmock"
	DBMysql         = "mysql"
	MQRabbit        = "rabbitmq"
	MQRabbitSecured = "rabbitmq_secured"
	CacheRedis      = "redis"
)

type Config struct {
	viper.Viper
	// ConfigName string
	// ConfigType string
	// ConfigPath string
	LogLevel int
	ApiPort  int
	ChatUrl  string

	HtmlTemplatesPath string

	DBtype          string
	DBhost          string
	DBport          int
	DBname          string
	DBuser          string
	DBpass          string
	DBMigrationPath string

	QueueType  string
	QueueHost  string
	QueuePort  int
	QueueUser  string
	QueuePass  string
	QueueVHost string

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
		LogLevel: -1,
		ApiPort:  80,
		ChatUrl:  "ws://127.0.0.1:8090/ws/chat",

		HtmlTemplatesPath: "html_tmpl",

		DBtype:          DBMysql,
		DBhost:          "localhost",
		DBport:          3306,
		DBname:          "sntest",
		DBuser:          "testuser",
		DBpass:          "testpass",
		DBMigrationPath: "migpath",

		QueueType:  MQRabbit,
		QueueHost:  "127.0.0.1",
		QueuePort:  5672,
		QueueUser:  "testuser",
		QueuePass:  "testpass",
		QueueVHost: "",

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

	_, ok := os.LookupEnv("PORT")
	if ok {
		_ = conf.BindEnv("apiport", "PORT")
	} else {
		_ = conf.BindEnv("apiport", "SOCIAL_NETWORK_APIPORT")
	}
	_ = conf.BindEnv("chat_url", "SOCIAL_NETWORK_CHAT_URL")

	_ = conf.BindEnv("html_tmpl_path", "SOCIAL_NETWORK_HTMLTMPLPATH")

	_ = conf.BindEnv("dbtype", "SOCIAL_NETWORK_DBTYPE")
	_ = conf.BindEnv("dbhost", "SOCIAL_NETWORK_DBHOST")
	_ = conf.BindEnv("dbport", "SOCIAL_NETWORK_DBPORT")
	_ = conf.BindEnv("dbname", "SOCIAL_NETWORK_DBNAME")
	_ = conf.BindEnv("dbuser", "SOCIAL_NETWORK_DBUSER")
	_ = conf.BindEnv("dbpass", "SOCIAL_NETWORK_DBPASS")
	_ = conf.BindEnv("dbmigrationpath", "SOCIAL_NETWORK_DBMIGRATIONPATH")

	_ = conf.BindEnv("queuetype", "SOCIAL_NETWORK_QUEUETYPE")
	_ = conf.BindEnv("queuehost", "SOCIAL_NETWORK_QUEUEHOST")
	_ = conf.BindEnv("queueport", "SOCIAL_NETWORK_QUEUEPORT")
	_ = conf.BindEnv("queueuser", "SOCIAL_NETWORK_QUEUEUSER")
	_ = conf.BindEnv("queuepass", "SOCIAL_NETWORK_QUEUEPASS")
	_ = conf.BindEnv("queuevhost", "SOCIAL_NETWORK_QUEUEVHOST")

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
	conf.SetDefault("chat_url", conf.ChatUrl)

	conf.SetDefault("html_tmpl_path", conf.HtmlTemplatesPath)

	conf.SetDefault("dbtype", conf.DBtype)
	conf.SetDefault("dbhost", conf.DBhost)
	conf.SetDefault("dbport", conf.DBport)
	conf.SetDefault("dbname", conf.DBname)
	conf.SetDefault("dbuser", conf.DBuser)
	conf.SetDefault("dbpass", conf.DBpass)
	conf.SetDefault("dbmigrationpath", conf.DBMigrationPath)

	conf.SetDefault("queuetype", conf.QueueType)
	conf.SetDefault("queuehost", conf.QueueHost)
	conf.SetDefault("queueport", conf.QueuePort)
	conf.SetDefault("queueuser", conf.QueueUser)
	conf.SetDefault("queuepass", conf.QueuePass)
	conf.SetDefault("queuevhost", conf.QueueVHost)

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
	conf.ChatUrl = conf.GetString("chat_url")

	conf.HtmlTemplatesPath = conf.GetString("html_tmpl_path")

	conf.DBtype = conf.GetString("dbtype")
	conf.DBhost = conf.GetString("dbhost")
	conf.DBport = conf.GetInt("dbport")
	conf.DBname = conf.GetString("dbname")
	conf.DBuser = conf.GetString("dbuser")
	conf.DBpass = conf.GetString("dbpass")
	conf.DBMigrationPath = conf.GetString("dbmigrationpath")

	conf.QueueType = conf.GetString("queuetype")
	conf.QueueHost = conf.GetString("queuehost")
	conf.QueuePort = conf.GetInt("queueport")
	conf.QueueUser = conf.GetString("queueuser")
	conf.QueuePass = conf.GetString("queuepass")
	conf.QueueVHost = conf.GetString("queuevhost")

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
