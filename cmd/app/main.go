package main

import (
	"log"
	"time"

	"github.com/avbar/mitemp/internal/config"
	"github.com/avbar/mitemp/internal/handler"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal("config init: ", err)
	}

	log.Print("creating handler")
	handler, err := handler.NewHandler(config.ConfigData.Sensors)
	for err != nil {
		log.Fatal("error creating handler: ", err)
	}
	log.Print("handler created")

	for {
		handler.Handle()
		time.Sleep(10 * time.Second)
	}
}
