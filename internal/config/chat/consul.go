package chat

func (conf *Config) ServerAddr() string {
	return conf.ConsulServerAddr
}

func (conf *Config) AgentAddr() string {
	return conf.ConsulAgentAddr
}

func (conf *Config) ServiceName() string {
	return conf.ConsulServiceName
}

func (conf *Config) ServiceID() string {
	return conf.ConsulServiceID
}
