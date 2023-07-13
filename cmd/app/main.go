package main

import (
	"log"
	"time"

	"github.com/avbar/mitemp/internal/handler"
)

var (
	sensor = handler.Sensor{
		Name: "My sensor",
		MAC:  "a4:c1:38:8a:d6:c2",
	}
)

func main() {
	log.Print("creating handler")
	handler, err := handler.NewHandler(sensor)
	for err != nil {
		log.Fatal("error creating handler: ", err)
	}
	log.Print("handler created")

	for {
		r, err := handler.GetReading()
		if err != nil {
			log.Print("error getting readings: ", err)
		} else {
			log.Printf("readings: %+v", r)
		}

		time.Sleep(10 * time.Second)
	}
}
