package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/tyghr/logger"
	"github.com/tyghr/social_network/internal/config"
	consul "github.com/tyghr/social_network/internal/consul/agent"
	"github.com/tyghr/social_network/internal/httpserver"
	"github.com/tyghr/social_network/internal/storage"
	"github.com/tyghr/social_network/internal/storage/cache/redis"
	"github.com/tyghr/social_network/internal/storage/db/mysql"
	"github.com/tyghr/social_network/internal/storage/queue/rabbitmq"
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

	db, err := mysql.OpenConn(conf, lgr)
	if err != nil {
		lgr.Fatal(err)
	}
	queue := rabbitmq.New(conf, lgr)
	cache := redis.New(conf.CacheNodes, conf.CacheClustered, conf.CachePass, lgr)
	stor := storage.New(db, queue, cache)

	// "append" queue processing
	chPostBuckets, err := queue.ReadPostAppendBuckets(context.TODO())
	if err != nil {
		lgr.Fatalw("ReadPostAppendBuckets", "error", err)
		return
	}
	go func() {
		lgr.Debugw("consumer postBuckets started")
		for pb := range chPostBuckets {
			err := cache.AddPostToSubscribers(context.TODO(), pb.Post, pb.Subscribers)
			if err != nil {
				lgr.Fatalw("AddPostToSubscribers", "error", err)
				break
			}
		}
		lgr.Debugw("consumer postBuckets ended")
	}()

	// "rebuild" queue processing
	chRebuildFeed, err := queue.ReadFeedRebuild(context.TODO())
	if err != nil {
		lgr.Fatalw("ReadFeedRebuild", "error", err)
		return
	}
	go func() {
		lgr.Debugw("consumer feedRebuildBuckets started")
		for sub := range chRebuildFeed {
			nr, posts, err := db.GetFeedRebuild(context.TODO(), sub)
			if err != nil {
				lgr.Fatalw("GetFeedRebuild", "error", err)
				break
			}
			if nr {
				err = cache.RebuildFeed(context.TODO(), sub, posts)
				if err != nil {
					lgr.Fatalw("RebuildFeed", "error", err)
					break
				}
			}
		}
		lgr.Debugw("consumer feedRebuildBuckets ended")
	}()

	server := httpserver.NewServer(stor, conf, lgr)
	lgr.Debug("start listening...")
	if err := http.ListenAndServe(
		fmt.Sprintf("0.0.0.0:%d", conf.ApiPort),
		server,
	); err != nil {
		lgr.Fatal(err)
	}
}
