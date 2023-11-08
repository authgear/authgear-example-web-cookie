package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/log"

	"github.com/oursky/authgear-exmaple-web-cookie/pkg/authgear"
	"github.com/oursky/authgear-exmaple-web-cookie/pkg/config"
	"github.com/oursky/authgear-exmaple-web-cookie/pkg/handler"
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

	router := httproute.NewRouter()
	router.Add(httproute.Route{
		Methods:     []string{"GET"},
		PathPattern: "/ping",
	}, &handler.PingHandler{})
	router.Add(httproute.Route{
		Methods:     []string{"GET"},
		PathPattern: "/",
	}, &handler.IndexHandler{
		AuthgearClient:   authgearClient,
		AuthgearEndpoint: envCfg.AuthgearEndpoint,
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
