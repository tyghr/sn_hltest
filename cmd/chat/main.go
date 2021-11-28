package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/tyghr/logger"
	config "github.com/tyghr/social_network/internal/config/chat"
	httpChat "github.com/tyghr/social_network/internal/httpserver/chat"
	consul "github.com/tyghr/social_network/internal/infra/consul/agent"
	redisChat "github.com/tyghr/social_network/internal/storage/cache/redis/chat"
)

func main() {
	if tz := os.Getenv("TZ"); tz != "" {
		var err error
		time.Local, err = time.LoadLocation(tz)
		if err != nil {
			panic(err)
		}
	}

	conf := config.NewConfig()
	if err := conf.ReadAllSettings(); err != nil {
		panic(err)
	}

	lgr := logger.NewLogger(conf.LogLevel, logger.ServiceLogger)

	// consul part
	consulClient, err := consul.NewClient(conf)
	if err != nil {
		lgr.Fatal(err)
	}

	if err = consulClient.Register(); err != nil {
		lgr.Fatal(err)
	}

	defer func() {
		if err = consulClient.Deregister(); err != nil {
			lgr.Fatal(err)
		}

		lgr.Debug("service auth deregister in consul")
	}()

	cache := redisChat.New(conf.CacheNodes, conf.CacheClustered, conf.CachePass, lgr)

	chat := httpChat.NewChatServer(cache, conf, lgr)
	lgr.Debug("start listening...")
	if err := http.ListenAndServe(
		fmt.Sprintf("0.0.0.0:%d", conf.ApiPort),
		chat,
	); err != nil {
		lgr.Fatal(err)
	}
}
