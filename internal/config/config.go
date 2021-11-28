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
	LogLevel          int
	ApiPort           int
	HtmlTemplatesPath string

	ChatUrl string

	DBConfig     *DBConfig
	QueueConfig  *QueueConfig
	CacheConfig  *CacheConfig
	ConsulConfig *ConsulConfig
	ZabbixConfig *ZabbixConfig
}

type DBConfig struct {
	Type          string
	Host          string
	Port          int
	Name          string
	User          string
	Pass          string
	MigrationPath string
}

type QueueConfig struct {
	Type  string
	Host  string
	Port  int
	User  string
	Pass  string
	VHost string
}

type CacheConfig struct {
	Type      string
	Nodes     []string
	Pass      string
	Clustered bool
}

type ConsulConfig struct {
	ServerAddr  string
	ServiceName string
	ServiceID   string
	AgentAddr   string
}

type ZabbixConfig struct {
	ServerHost string
	Port       int
	HostName   string
}

func NewConfig() *Config {
	return &Config{
		Viper: *viper.New(),
		// ConfigName: "config",
		// ConfigType: "json",
		// ConfigPath: "./", // change in prod to "/etc/social_network/"
		LogLevel:          -1,
		ApiPort:           80,
		HtmlTemplatesPath: "html_tmpl",

		ChatUrl: "ws://127.0.0.1:8090/ws/chat",

		DBConfig: &DBConfig{
			Type:          DBMysql,
			Host:          "localhost",
			Port:          3306,
			Name:          "sntest",
			User:          "testuser",
			Pass:          "testpass",
			MigrationPath: "migpath",
		},
		QueueConfig: &QueueConfig{
			Type:  MQRabbit,
			Host:  "127.0.0.1",
			Port:  5672,
			User:  "testuser",
			Pass:  "testpass",
			VHost: "",
		},
		CacheConfig: &CacheConfig{
			Type:      CacheRedis,
			Nodes:     []string{"redis_node_0:6379", "redis_node_1:6379", "redis_node_2:6379", "redis_node_3:6379", "redis_node_4:6379", "redis_node_5:6379"},
			Pass:      "testpass",
			Clustered: false,
		},
		ConsulConfig: &ConsulConfig{
			ServerAddr:  "",
			ServiceName: "",
			ServiceID:   "",
			AgentAddr:   "",
		},
		ZabbixConfig: &ZabbixConfig{
			ServerHost: "",
			Port:       0,
			HostName:   "",
		},
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

	_ = conf.BindEnv("zabbix_server_host", "ZABBIX_SERVER_HOST")
	_ = conf.BindEnv("zabbix_port", "ZABBIX_PORT")
	_ = conf.BindEnv("zabbix_host_name", "ZABBIX_HOST_NAME")
}

func (conf *Config) setDefaults() {
	conf.SetDefault("apiport", conf.ApiPort)
	conf.SetDefault("loglevel", conf.LogLevel)
	conf.SetDefault("chat_url", conf.ChatUrl)

	conf.SetDefault("html_tmpl_path", conf.HtmlTemplatesPath)

	conf.SetDefault("dbtype", conf.DBConfig.Type)
	conf.SetDefault("dbhost", conf.DBConfig.Host)
	conf.SetDefault("dbport", conf.DBConfig.Port)
	conf.SetDefault("dbname", conf.DBConfig.Name)
	conf.SetDefault("dbuser", conf.DBConfig.User)
	conf.SetDefault("dbpass", conf.DBConfig.Pass)
	conf.SetDefault("dbmigrationpath", conf.DBConfig.MigrationPath)

	conf.SetDefault("queuetype", conf.QueueConfig.Type)
	conf.SetDefault("queuehost", conf.QueueConfig.Host)
	conf.SetDefault("queueport", conf.QueueConfig.Port)
	conf.SetDefault("queueuser", conf.QueueConfig.User)
	conf.SetDefault("queuepass", conf.QueueConfig.Pass)
	conf.SetDefault("queuevhost", conf.QueueConfig.VHost)

	conf.SetDefault("cachetype", conf.CacheConfig.Type)
	conf.SetDefault("cachenodes", conf.CacheConfig.Nodes)
	conf.SetDefault("cachepass", conf.CacheConfig.Pass)
	conf.SetDefault("cacheclustered", conf.CacheConfig.Clustered)

	conf.SetDefault("consul_server_addr", conf.ConsulConfig.ServerAddr)
	conf.SetDefault("consul_service_name", conf.ConsulConfig.ServiceName)
	conf.SetDefault("consul_service_id", conf.ConsulConfig.ServiceID)
	conf.SetDefault("consul_agent_addr", conf.ConsulConfig.AgentAddr)

	conf.SetDefault("zabbix_server_host", conf.ZabbixConfig.ServerHost)
	conf.SetDefault("zabbix_port", conf.ZabbixConfig.Port)
	conf.SetDefault("zabbix_host_name", conf.ZabbixConfig.HostName)
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

	conf.DBConfig.Type = conf.GetString("dbtype")
	conf.DBConfig.Host = conf.GetString("dbhost")
	conf.DBConfig.Port = conf.GetInt("dbport")
	conf.DBConfig.Name = conf.GetString("dbname")
	conf.DBConfig.User = conf.GetString("dbuser")
	conf.DBConfig.Pass = conf.GetString("dbpass")
	conf.DBConfig.MigrationPath = conf.GetString("dbmigrationpath")

	conf.QueueConfig.Type = conf.GetString("queuetype")
	conf.QueueConfig.Host = conf.GetString("queuehost")
	conf.QueueConfig.Port = conf.GetInt("queueport")
	conf.QueueConfig.User = conf.GetString("queueuser")
	conf.QueueConfig.Pass = conf.GetString("queuepass")
	conf.QueueConfig.VHost = conf.GetString("queuevhost")

	conf.CacheConfig.Type = conf.GetString("cachetype")
	conf.CacheConfig.Nodes = conf.GetStringSlice("cachenodes")
	conf.CacheConfig.Pass = conf.GetString("cachepass")
	conf.CacheConfig.Clustered = conf.GetBool("cacheclustered")

	conf.ConsulConfig.ServerAddr = conf.GetString("consul_server_addr")
	conf.ConsulConfig.ServiceName = conf.GetString("consul_service_name")
	conf.ConsulConfig.ServiceID = conf.GetString("consul_service_id")
	conf.ConsulConfig.AgentAddr = conf.GetString("consul_agent_addr")

	conf.ZabbixConfig.ServerHost = conf.GetString("zabbix_server_host")
	conf.ZabbixConfig.Port = conf.GetInt("zabbix_port")
	conf.ZabbixConfig.HostName = conf.GetString("zabbix_host_name")

	return nil
}
