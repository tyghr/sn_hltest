package config

func (conf *Config) ServerAddr() string {
	return conf.ConsulConfig.ServerAddr
}

func (conf *Config) AgentAddr() string {
	return conf.ConsulConfig.AgentAddr
}

func (conf *Config) ServiceName() string {
	return conf.ConsulConfig.ServiceName
}

func (conf *Config) ServiceID() string {
	return conf.ConsulConfig.ServiceID
}
