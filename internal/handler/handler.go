package handler

import (
	"context"
	"encoding/binary"
	"log"
	"strings"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
)

var (
	// UUID of sensor characteristic with temperature/humidity/voltage
	tempUUID = ble.MustParse("ebe0ccc1-7a0a-4b0c-8a1a-6ff2997da3a6")
)

type Handler struct {
	sensor Sensor
}

func NewHandler(sensor Sensor) (*Handler, error) {
	h := &Handler{
		sensor: sensor,
	}

	d, err := linux.NewDevice()
	for err != nil {
		return nil, err
	}
	ble.SetDefaultDevice(d)

	return h, nil
}

func (h *Handler) handleReading(ch chan<- Reading) func([]byte) {
	return func(b []byte) {
		// Bytes 0-1 are temperature, byte 2 is humidity, bytes 3-4 are voltage
		r := Reading{
			Temperature: float64(binary.LittleEndian.Uint16(b[0:2])) / 100,
			Humidity:    float64(b[2]),
			Voltage:     float64(binary.LittleEndian.Uint16(b[3:5])) / 1000,
		}
		ch <- r
	}
}

func (h *Handler) GetReading() (Reading, error) {
	log.Printf("getting readings for %q", h.sensor.Name)

	var r Reading

	// Filter sensor by MAC address
	filter := func(a ble.Advertisement) bool {
		return strings.EqualFold(a.Addr().String(), h.sensor.MAC)
	}

	// Scan for specified duration, or until interrupted by user
	log.Print("connecting...")
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), 20*time.Second))
	cln, err := ble.Connect(ctx, filter)
	if err != nil {
		return r, err
	}
	log.Printf("connected to %q", h.sensor.Name)

	done := make(chan struct{})

	// Normally, the connection is disconnected by us after our exploration.
	// However, it can be asynchronously disconnected by the remote peripheral.
	// So we wait(detect) the disconnection in the go routine.
	go func() {
		<-cln.Disconnected()
		log.Printf("disconnected from %q", h.sensor.Name)
		close(done)
	}()

	log.Print("discovering profile...")
	p, err := cln.DiscoverProfile(true)
	if err != nil {
		return r, err
	}

	log.Print("finding characteristic...")
	c := p.FindCharacteristic(ble.NewCharacteristic(tempUUID))

	log.Print("reading data...")
	ch := make(chan Reading)
	err = cln.Subscribe(c, false, h.handleReading(ch))
	if err != nil {
		return r, err
	}

	select {
	case r = <-ch:
	case <-done:
		return r, nil
	}

	log.Print("disconnecting...")
	cln.CancelConnection()
	<-done

	return r, nil
}
