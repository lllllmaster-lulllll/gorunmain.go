package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	version = "2.10.5"
	dvs = []Device{
		{
			ID:       1,
			Mac:      "00:00:00:00:00:01",
			Firmware: "3.1.1",
		},
		{
			ID:       2,
			Mac:      "00:00:00:00:00:02",
			Firmware: "3.1.2",
		},
		{
			ID:       3,
			Mac:      "00:00:00:00:00:03",
			Firmware: "3.1.3",
		},
		{
			ID:       4,
			Mac:      "00:00:00:00:00:04",
			Firmware: "3.1.4",
		},
	}

}

var dvs []Device
var version string

type Device struct {
	ID int `json:"id"`

	Mac string `json:"mac"`

	Firmware string `json:"firmware"`
}

type metrics struct {
	devices prometheus.Gauge
	info    *prometheus.GaugeVec
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		devices: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "myapp",
			Name:      "devices",
			Help:      "Number of devices.",
		}),
		info: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "myapp",
			Name:      "info",
			Help:      "Information about the My App environment.",
		},
			[]string{"version", "version02", "version03", "version04"},
		),
	}
	reg.MustRegister(m.devices, m.info)

	return m
}

func main() {
	reg := prometheus.NewRegistry()
	m := NewMetrics(reg)
	m.devices.Set(float64(len(dvs)))
	version02 := "3.1.1"
	version03 := "3.1.2"
	version04 := "3.1.3"
	m.info.With(prometheus.Labels{"version": version, "version02": version02, "version03": version03, "version04": version04}).Set(1)
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	// http.Handle("/metrics", promHandler)
	// http.HandleFunc("/devices", getDevices)
	// http.ListenAndServe(":9000", nil)
	dMux := http.NewServeMux()
	dMux.HandleFunc("/devices", getDevices)

	pMux := http.NewServeMux()
	pMux.Handle("/metrics", promHandler)

	go func() {
		log.Fatal(http.ListenAndServe(":9000", pMux))
	}()

	go func() {
		log.Fatal(http.ListenAndServe(":9001", dMux))
	}()

	select {}
}
func getDevices(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(dvs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
