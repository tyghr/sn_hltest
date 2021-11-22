package config

func (conf *Config) ServerAddr() string {
	return conf.ConsulServerAddr
}

func (conf *Config) AgentAddr() string {
	return ""
}

func (conf *Config) ServiceName() string {
	return conf.ConsulServiceName
}

func (conf *Config) ServiceID() string {
	return ""
}
