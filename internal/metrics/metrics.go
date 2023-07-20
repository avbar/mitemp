package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	GaugeTemperature = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "app",
		Subsystem: "sensors",
		Name:      "gauge_temperature_degrees",
	},
		[]string{"sensor"},
	)
	GaugeHumidity = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "app",
		Subsystem: "sensors",
		Name:      "gauge_humidity_percents",
	},
		[]string{"sensor"},
	)
	GaugeVoltage = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "app",
		Subsystem: "sensors",
		Name:      "gauge_voltage_volts",
	},
		[]string{"sensor"},
	)
	GaugeBattery = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "app",
		Subsystem: "sensors",
		Name:      "gauge_battery_percents",
	},
		[]string{"sensor"},
	)
)

func New() http.Handler {
	return promhttp.Handler()
}
