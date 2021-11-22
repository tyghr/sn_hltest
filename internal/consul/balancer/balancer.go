package balancer

import (
	"strconv"
	"time"

	config "github.com/tyghr/social_network/internal/config/consul"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/tyghr/logger"
)

type Client struct {
	client      *consulapi.Client
	serviceName string
	srvList     *ServerAvailableList

	lgr logger.Logger

	ticker *time.Ticker
}

func NewClient(cfg config.ConsulConfig, lgr logger.Logger) (*Client, error) {
	consulCfg := consulapi.DefaultConfig()
	consulCfg.Address = cfg.ServerAddr()

	client, err := consulapi.NewClient(consulCfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		client:      client,
		serviceName: cfg.ServiceName(),
		srvList:     NewServerAvailableList(),
		lgr:         lgr,
		ticker:      time.NewTicker(time.Second * 10),
	}, nil
}

func (c *Client) HealthCheck() {
	for range c.ticker.C {
		health, _, err := c.client.Health().Service(c.serviceName, "", false, nil)
		if err != nil {
			c.lgr.Error("cannot to observe available services via Consul", err)

			continue
		}

		var servers []string
		for _, item := range health {
			addr := item.Service.Address + ":" + strconv.Itoa(item.Service.Port)
			servers = append(servers, addr)
		}

		c.srvList.Update(servers)
	}
}

func (c *Client) Stop() {
	c.ticker.Stop()
}
