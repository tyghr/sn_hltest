package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/tyghr/logger"
	"github.com/tyghr/social_network/internal/config"
	"github.com/tyghr/social_network/internal/httpserver"
	"github.com/tyghr/social_network/internal/storage/db/mysql"
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

	lgr := logger.NewLogger(conf.LogLevel, conf.LogLevel == -1)

	db, err := mysql.OpenConn(conf, lgr)
	if err != nil {
		lgr.Fatal(err)
	}

	server := httpserver.NewServer(db, conf, lgr)
	lgr.Debug("start listening...")
	if err := http.ListenAndServe(
		fmt.Sprintf("0.0.0.0:%d", conf.ApiPort),
		server,
	); err != nil {
		lgr.Fatal(err)
	}
}
