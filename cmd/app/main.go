package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/avbar/mitemp/internal/config"
	"github.com/avbar/mitemp/internal/handler"
	"github.com/avbar/mitemp/internal/logger"
	"github.com/avbar/mitemp/internal/metrics"
	"go.uber.org/zap"
)

var develMode = flag.Bool("devel", true, "development mode")

func main() {
	flag.Parse()

	logger.Init(*develMode)

	err := config.Init()
	if err != nil {
		logger.Fatal("config init error", zap.Error(err))
	}

	handler, err := handler.NewHandler(config.ConfigData.Sensors)
	for err != nil {
		logger.Fatal("error creating handler", zap.Error(err))
	}

	go func() {
		http.Handle("/metrics", metrics.New())
		http.ListenAndServe(fmt.Sprintf(":%d", config.ConfigData.Port), nil)
	}()

	handler.Handle()
}
