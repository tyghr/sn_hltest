package zabbix

import (
	"context"
	"fmt"
	"strconv"
	"time"

	zbx "github.com/blacked/go-zabbix"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/tyghr/logger"
	"github.com/tyghr/social_network/internal/config"
)

type Client struct {
	sender   *zbx.Sender
	hostName string

	lgr logger.Logger
}

func NewClient(cfg *config.ZabbixConfig, logger logger.Logger) *Client {
	return &Client{
		sender:   zbx.NewSender(cfg.ServerHost, cfg.Port),
		hostName: cfg.HostName,
		lgr:      logger,
	}
}

func (c *Client) Publish(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			pkg := c.collectMetrics()

			if _, err := c.sender.Send(pkg); err != nil {
				c.lgr.Error("send pkg to zabbix", err)
			}
		}
	}

}

func (c *Client) collectMetrics() *zbx.Packet {
	metrics := make([]*zbx.Metric, 0, 1)

	memMetrics, err := c.collectMemoryMetrics()
	if err != nil {
		c.lgr.Error("get MEMORY info", err)
	}
	metrics = append(metrics, memMetrics...)

	cpuMetrics, err := c.collectCPUMetrics()
	if err != nil {
		c.lgr.Error("get CPU info", err)
	}
	metrics = append(metrics, cpuMetrics...)

	return zbx.NewPacket(metrics)
}

func (c *Client) collectMemoryMetrics() ([]*zbx.Metric, error) {
	metrics := make([]*zbx.Metric, 0, 4)

	mem, err := memory.Get()
	if err != nil {
		return nil, err
	}

	metrics = append(metrics, zbx.NewMetric(c.hostName, "sn_server_mem_total", strconv.FormatUint(mem.Total, 10), time.Now().Unix()))
	metrics = append(metrics, zbx.NewMetric(c.hostName, "sn_server_mem_used", strconv.FormatUint(mem.Used, 10), time.Now().Unix()))
	metrics = append(metrics, zbx.NewMetric(c.hostName, "sn_server_mem_free", strconv.FormatUint(mem.Free, 10), time.Now().Unix()))

	return metrics, nil
}

func (c *Client) collectCPUMetrics() ([]*zbx.Metric, error) {
	metrics := make([]*zbx.Metric, 0, 4)

	cpuBefore, err := cpu.Get()
	if err != nil {
		return nil, err
		//c.logger.Error("get CPU info", err)
	}

	time.Sleep(time.Second * 1)

	cpuAfter, err := cpu.Get()
	if err != nil {
		return nil, err
	}

	total := float64(cpuAfter.Total - cpuBefore.Total)
	cpuUser := float64(cpuAfter.User-cpuBefore.User) / total * 100
	cpuSystem := float64(cpuAfter.System-cpuBefore.System) / total * 100
	cpuIdle := float64(cpuAfter.Idle-cpuBefore.Idle) / total * 100

	metrics = append(metrics, zbx.NewMetric(c.hostName, "sn_server_cpu_user", fmt.Sprintf("%f", cpuUser), time.Now().Unix()))
	metrics = append(metrics, zbx.NewMetric(c.hostName, "sn_server_cpu_system", fmt.Sprintf("%f", cpuSystem), time.Now().Unix()))
	metrics = append(metrics, zbx.NewMetric(c.hostName, "sn_server_cpu_idle", fmt.Sprintf("%f", cpuIdle), time.Now().Unix()))

	return metrics, nil
}
