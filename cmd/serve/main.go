package main

import (
	"errors"
	"net/http"
	"net/url"
	"os"

	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/log"

	"github.com/authgear/authgear-exmaple-web-cookie/pkg/authgear"
	"github.com/authgear/authgear-exmaple-web-cookie/pkg/config"
	"github.com/authgear/authgear-exmaple-web-cookie/pkg/handler"
)

func main() {
	envCfg, err := config.LoadConfigFromEnv()
	if err != nil {
		panic(err)
	}

	loggerFactory := log.NewFactory(
		envCfg.LogLevel,
	)

	mainLogger := loggerFactory.New("main")

	authgearClient := authgear.NewClient(envCfg.AuthgearEndpoint)
	authgearEndpointURL, err := url.Parse(envCfg.AuthgearEndpoint)
	if err != nil {
		panic("invalid endpoint")
	}

	router := httproute.NewRouter()
	router.Add(httproute.Route{
		Methods:     []string{"GET", "POST"},
		PathPattern: "/",
	}, &handler.IndexHandler{
		Logger:           loggerFactory.New("index"),
		AuthgearClient:   authgearClient,
		AuthgearEndpoint: authgearEndpointURL,
		DefaultClientID:  envCfg.DefaultClientID,
	})

	mainLogger.Info("starting server on ", envCfg.ListenAddr)
	err = http.ListenAndServe(envCfg.ListenAddr, router)
	if errors.Is(err, http.ErrServerClosed) {
		mainLogger.Info("server closed")
	} else if err != nil {
		mainLogger.WithError(err).Error("error")
		os.Exit(1)
	}
}
