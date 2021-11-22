package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/tyghr/logger"
	config "github.com/tyghr/social_network/internal/config/counters"
	consul "github.com/tyghr/social_network/internal/consul/agent"
	httpCounters "github.com/tyghr/social_network/internal/httpserver/counters"
	"github.com/tyghr/social_network/internal/model"
	"github.com/tyghr/social_network/internal/storage"
	redis "github.com/tyghr/social_network/internal/storage/cache/redis/counters"
	rmqCounters "github.com/tyghr/social_network/internal/storage/queue/rabbitmq/counters"
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

	queue := rmqCounters.New(conf, lgr)
	cache := redis.New(conf.CacheNodes, conf.CacheClustered, conf.CachePass, lgr)
	stor := storage.NewCounters(queue, cache)

	// queue consumer
	chCounters, err := queue.ReadCountersBuckets(context.TODO())
	if err != nil {
		lgr.Fatalw("ReadFeedRebuild", "error", err)
		return
	}
	go func() {
		lgr.Debugw("consumer countersBuckets started")
		for cmd := range chCounters {
			ctx := context.TODO()
			switch cmd.Command {
			case model.CounterCmdTotalInc:
				err := cache.IncCounters(ctx, model.CounterCmdTotalInc, cmd.Subscribers)
				if err != nil {
					lgr.Fatalw("IncCounters (Total)", "error", err)
					break
				}
			case model.CounterCmdCursorInc:
				err := cache.IncCounters(ctx, model.CounterCmdCursorInc, cmd.Subscribers)
				if err != nil {
					lgr.Fatalw("IncCounters (Cursor)", "error", err)
					break
				}
			case model.CounterCmdCursorUpdate:
				err := cache.UpdateCursorCounter(ctx, cmd.Subscribers)
				if err != nil {
					lgr.Fatalw("UpdateCursorCounter", "error", err)
					break
				}
			}
		}
		lgr.Debugw("consumer countersBuckets ended")
	}()

	srvCounters := httpCounters.NewCountersServer(stor, conf, lgr)
	lgr.Debug("start listening...")
	if err := http.ListenAndServe(
		fmt.Sprintf("0.0.0.0:%d", conf.ApiPort),
		srvCounters,
	); err != nil {
		lgr.Fatal(err)
	}
}
