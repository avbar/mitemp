package main

import (
	"flag"

	"github.com/avbar/mitemp/internal/config"
	"github.com/avbar/mitemp/internal/handler"
	"github.com/avbar/mitemp/internal/logger"
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

	handler.Handle()
}
