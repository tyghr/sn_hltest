package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	TestMock = "testmock"
	DBMysql  = "mysql"
)

type Config struct {
	viper.Viper
	// ConfigName string
	// ConfigType string
	// ConfigPath string
	LogLevel int
	ApiPort  int
	DBtype   string
	DBhost   string
	DBport   int
	DBname   string
	DBuser   string
	DBpass   string
}

func NewConfig() *Config {
	return &Config{
		Viper: *viper.New(),
		// ConfigName: "config",
		// ConfigType: "json",
		// ConfigPath: "./", // change in prod to "/etc/social_network/"
		LogLevel: -1,
		ApiPort:  80,
		DBtype:   DBMysql,
		DBhost:   "localhost",
		DBport:   3306,
		DBname:   "sntest",
		DBuser:   "root",
		DBpass:   "secretpass",
	}
}

func init() {
	// pflag.String("config", "", "config file path")
	pflag.Int("loglevel", 0, "debug level (debug:-1 .. fatal:5)")
	pflag.Int("apiport", 0, "API port")
	pflag.String("dbtype", "", "DB type")
	pflag.String("dbhost", "", "DB host")
	pflag.String("dbport", "", "DB port")
	pflag.String("dbname", "", "DB name")
	pflag.String("dbuser", "", "DB user")
	pflag.String("dbpass", "", "DB pass")
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
}

func (conf *Config) bindAllFlags() {
	pflag.Parse()
	// _ = conf.BindPFlag("config", pflag.Lookup("config"))
	_ = conf.BindPFlag("apiport", pflag.Lookup("apiport"))
	_ = conf.BindPFlag("loglevel", pflag.Lookup("loglevel"))
	_ = conf.BindPFlag("dbtype", pflag.Lookup("dbtype"))
	_ = conf.BindPFlag("dbhost", pflag.Lookup("dbhost"))
	_ = conf.BindPFlag("dbport", pflag.Lookup("dbport"))
	_ = conf.BindPFlag("dbname", pflag.Lookup("dbname"))
	_ = conf.BindPFlag("dbuser", pflag.Lookup("dbuser"))
	_ = conf.BindPFlag("dbpass", pflag.Lookup("dbpass"))
}

func (conf *Config) setDefaults() {
	// conf.SetDefault("config", fmt.Sprintf("%s/%s.%s", strings.TrimSuffix(conf.ConfigPath, "/"), conf.ConfigName, conf.ConfigType))
	conf.SetDefault("apiport", conf.ApiPort)
	conf.SetDefault("loglevel", conf.LogLevel)
	conf.SetDefault("dbtype", conf.DBtype)
	conf.SetDefault("dbhost", conf.DBhost)
	conf.SetDefault("dbport", conf.DBport)
	conf.SetDefault("dbname", conf.DBname)
	conf.SetDefault("dbuser", conf.DBuser)
	conf.SetDefault("dbpass", conf.DBpass)
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
	conf.bindAllFlags()

	// flagConfig := conf.GetString("config")
	// if flagConfig != "" {
	// 	conf.ConfigPath = path.Dir(flagConfig)
	// 	conf.ConfigType = strings.TrimPrefix(path.Ext(flagConfig), ".")
	// 	if conf.ConfigType != "" {
	// 		conf.ConfigName = strings.TrimSuffix(path.Base(flagConfig), "."+conf.ConfigType)
	// 	} else {
	// 		conf.ConfigName = path.Base(flagConfig)
	// 	}
	// }

	// conf.SetConfigName(conf.ConfigName) // read config
	// conf.SetConfigType(conf.ConfigType)
	// conf.AddConfigPath(conf.ConfigPath)
	// if err := conf.ReadInConfig(); err != nil {
	// 	if errW := conf.WriteConfigAs(fmt.Sprintf("%s/%s.%s", conf.ConfigPath, conf.ConfigName, conf.ConfigType)); errW != nil {
	// 		return err
	// 	}
	// }

	conf.ApiPort = conf.GetInt("apiport")
	conf.LogLevel = conf.GetInt("loglevel")
	conf.DBtype = conf.GetString("dbtype")
	conf.DBhost = conf.GetString("dbhost")
	conf.DBport = conf.GetInt("dbport")
	conf.DBname = conf.GetString("dbname")
	conf.DBuser = conf.GetString("dbuser")
	conf.DBpass = conf.GetString("dbpass")

	return nil
}
