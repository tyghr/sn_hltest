package config

import (
	"github.com/spf13/viper"
)

const (
	TestMock = "testmock"
	DBMysql  = "mysql"
	MQRabbit = "rabbitmq"
)

type Config struct {
	viper.Viper
	// ConfigName string
	// ConfigType string
	// ConfigPath string
	LogLevel        int
	ApiPort         int
	DBtype          string
	DBhost          string
	DBport          int
	DBname          string
	DBuser          string
	DBpass          string
	DBMigrationPath string
	QueueType       string
	QueueHost       string
	QueuePort       int
	QueueUser       string
	QueuePass       string

	CacheNodes []string
	CachePass  string
}

func NewConfig() *Config {
	return &Config{
		Viper: *viper.New(),
		// ConfigName: "config",
		// ConfigType: "json",
		// ConfigPath: "./", // change in prod to "/etc/social_network/"
		LogLevel:        -1,
		ApiPort:         80,
		DBtype:          DBMysql,
		DBhost:          "localhost",
		DBport:          3306,
		DBname:          "sntest",
		DBuser:          "testuser",
		DBpass:          "testpass",
		DBMigrationPath: "/opt/snserver/migrations/mysql",

		QueueType: MQRabbit,
		QueueHost: "127.0.0.1",
		QueuePort: 5672,
		QueueUser: "testuser",
		QueuePass: "testpass",

		CacheNodes: []string{"redis_node_0:6379", "redis_node_1:6379", "redis_node_2:6379", "redis_node_3:6379", "redis_node_4:6379", "redis_node_5:6379"},
		CachePass:  "testpass",
	}
}

func (conf *Config) bindAllEnv() {
	_ = conf.BindEnv("loglevel", "SOCIAL_NETWORK_LOGLEVEL")
	_ = conf.BindEnv("apiport", "SOCIAL_NETWORK_APIPORT")
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

	_ = conf.BindEnv("cachenodes", "SOCIAL_NETWORK_CACHENODES")
	_ = conf.BindEnv("cachepass", "SOCIAL_NETWORK_CACHEPASS")
}

func (conf *Config) setDefaults() {
	conf.SetDefault("apiport", conf.ApiPort)
	conf.SetDefault("loglevel", conf.LogLevel)
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

	conf.SetDefault("cachenodes", conf.CacheNodes)
	conf.SetDefault("cachepass", conf.CachePass)
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

	conf.CacheNodes = conf.GetStringSlice("cachenodes")
	conf.CachePass = conf.GetString("cachepass")

	return nil
}
