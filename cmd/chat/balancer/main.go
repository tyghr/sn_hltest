package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/tyghr/logger"
	config "github.com/tyghr/social_network/internal/config/balancer"
	consul "github.com/tyghr/social_network/internal/consul/balancer"
)

type Backend struct {
	URL          *url.URL
	ReverseProxy *httputil.ReverseProxy
}

type Balancer struct {
	Backends map[string]*Backend
}

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
	consulClient, err := consul.NewClient(conf, lgr)
	if err != nil {
		lgr.Fatal(err)
	}
	go consulClient.HealthCheck()
	defer consulClient.Stop()

	// нужен прокси
	bl := Balancer{}

	// parse servers
	tokens := consulClient.SafeGetServerList()
	for _, tok := range tokens {
		serverUrl, err := url.Parse(tok)
		if err != nil {
			lgr.Fatal(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(serverUrl)
		proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
			lgr.Errorf("[%s] %s\n", serverUrl.Host, e.Error())
		}

		bl.Backends[tok] = &Backend{
			URL:          serverUrl,
			ReverseProxy: proxy,
		}

		lgr.Debugf("Configured server: %s\n", serverUrl)
	}

	// create http server
	server := http.Server{
		Addr: fmt.Sprintf(":%d", conf.ApiPort),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			peer, err := consulClient.GetAddr()
			if err != nil {
				http.Error(w, "Service not available", http.StatusServiceUnavailable)
			}
			bl.Backends[peer].ReverseProxy.ServeHTTP(w, r)
		}),
	}

	lgr.Debugf("Load Balancer started at :%d\n", conf.ApiPort)
	if err := server.ListenAndServe(); err != nil {
		lgr.Fatal(err)
	}
}
