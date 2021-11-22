package balancer

import (
	"fmt"
	"math/rand"
	"sync"
)

type ServerAvailableList struct {
	servers []string

	sync.Mutex
}

func NewServerAvailableList() *ServerAvailableList {
	return &ServerAvailableList{servers: make([]string, 0, 1)}
}

func (s *ServerAvailableList) GetAddr() (string, error) {
	s.Lock()
	defer s.Unlock()

	if len(s.servers) == 0 {
		return "", fmt.Errorf("all severs are not available")
	}

	return s.servers[rand.Intn(len(s.servers))], nil
}

func (s *ServerAvailableList) Update(servers []string) {
	s.Lock()
	defer s.Unlock()

	s.servers = servers
}

func (c *Client) SafeGetServerList() []string {
	c.srvList.Lock()
	defer c.srvList.Unlock()

	s := []string{}
	copy(s, c.srvList.servers)
	return s
}

func (c *Client) GetAddr() (string, error) {
	return c.srvList.GetAddr()
}
