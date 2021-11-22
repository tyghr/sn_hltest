package config

type ConsulConfig interface {
	ServerAddr() string
	AgentAddr() string
	ServiceName() string
	ServiceID() string
}
