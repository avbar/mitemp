package handler

import (
	"context"
	"encoding/binary"
	"strings"
	"time"

	"github.com/avbar/mitemp/internal/logger"
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
	"go.uber.org/zap"
)

var (
	// UUID of sensor characteristic with temperature/humidity/voltage
	tempUUID = ble.MustParse("ebe0ccc1-7a0a-4b0c-8a1a-6ff2997da3a6")

	readingTimeout = 15 * time.Second
	readingPause   = 10 * time.Second
)

type Handler struct {
	sensors []Sensor
}

func NewHandler(sensors []Sensor) (*Handler, error) {
	h := &Handler{
		sensors: sensors,
	}

	logger.Info("creating device")
	d, err := linux.NewDevice()
	for err != nil {
		return nil, err
	}
	logger.Info("device created")
	ble.SetDefaultDevice(d)

	return h, nil
}

func (h *Handler) handleReading(ch chan<- Reading) func([]byte) {
	return func(b []byte) {
		// Bytes 0-1 are temperature, byte 2 is humidity, bytes 3-4 are voltage
		r := Reading{
			Temperature: float32(binary.LittleEndian.Uint16(b[0:2])) / 100,
			Humidity:    uint8(b[2]),
			Voltage:     float32(binary.LittleEndian.Uint16(b[3:5])) / 1000,
		}

		// Battery level: 3.1V is 100%, 2.1V is 0%
		battery := (r.Voltage - 2.1) * 100
		if battery < 0 {
			battery = 0
		} else if battery > 100 {
			battery = 100
		}
		r.Battery = uint8(battery)

		ch <- r
	}
}

func (h *Handler) GetReading(sensor Sensor) (Reading, error) {
	logger.Info("getting readings", zap.String("sensor", sensor.Name))

	var r Reading

	// Filter sensor by MAC address
	filter := func(a ble.Advertisement) bool {
		return strings.EqualFold(a.Addr().String(), sensor.MAC)
	}

	// Scan for specified duration, or until interrupted by user
	logger.Debug("connecting...", zap.String("sensor", sensor.Name))
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), readingTimeout))
	cln, err := ble.Connect(ctx, filter)
	if err != nil {
		return r, err
	}
	logger.Debug("connected", zap.String("sensor", sensor.Name))

	done := make(chan struct{})

	// Normally, the connection is disconnected by us after our exploration.
	// However, it can be asynchronously disconnected by the remote peripheral.
	// So we wait(detect) the disconnection in the go routine.
	go func() {
		<-cln.Disconnected()
		logger.Debug("disconnected", zap.String("sensor", sensor.Name))
		close(done)
	}()

	logger.Debug("discovering profile...", zap.String("sensor", sensor.Name))
	p, err := cln.DiscoverProfile(true)
	if err != nil {
		return r, err
	}

	logger.Debug("finding characteristic...", zap.String("sensor", sensor.Name))
	c := p.FindCharacteristic(ble.NewCharacteristic(tempUUID))

	logger.Debug("reading data...", zap.String("sensor", sensor.Name))
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

	logger.Debug("disconnecting...", zap.String("sensor", sensor.Name))
	cln.CancelConnection()
	<-done

	return r, nil
}

func (h *Handler) Handle() {
	for {
		for _, sensor := range h.sensors {
			r, err := h.GetReading(sensor)
			if err != nil {
				logger.Error("error getting readings", zap.Error(err), zap.String("sensor", sensor.Name))
			} else {
				logger.Info("readings were taken", zap.String("sensor", sensor.Name), zap.Any("reading", r))
			}
		}

		time.Sleep(readingPause)
	}
}
