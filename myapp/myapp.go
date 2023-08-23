package myapp

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cast"
)

type metrics struct {
	thirdDownService *prometheus.GaugeVec
	thirdLiveService *prometheus.GaugeVec
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		thirdDownService: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "polaris_iaas",
			Name:      "down_services",
			Help:      "瀚海云第三方宕机服务",
		},
			[]string{"third_down_services", "third_down_services_name", "instance"},
		),
		thirdLiveService: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "polaris_iaas",
			Name:      "live_services",
			Help:      "瀚海云第三方正常服务",
		},
			[]string{"third_live_services"},
		),
	}
	reg.MustRegister(m.thirdDownService, m.thirdLiveService)

	return m
}

func MyApp() {
	reg := prometheus.NewRegistry()
	m := NewMetrics(reg)
	addrPort := make(map[string]string)

	// 添加需要探测的服务
	addrPort["local_redis"] = "192.168.1.15:6379"
	addrPort["local_prometheus"] = "192.168.1.15:9000"

	thirdDownService := make([]string, 0)
	thirdLiveService := make([]string, 0)
	thirdDownServiceName := make([]string, 0)

	for index, api := range addrPort {
		strList := strings.Split(api, ":")
		apilen := len(strList)
		if apilen != 2 {
			log.Println("api format error, api:", api)
			continue
		}
		addr := strList[0]
		port, err := cast.ToIntE(strList[1])
		if err != nil {
			log.Println(err)
			continue
		}

		live, err := DetectService(addr, port, time.Second)
		if err != nil {
			log.Println(err)
		}
		if !live {
			thirdDownServiceName = append(thirdDownServiceName, index)
			thirdDownService = append(thirdDownService, fmt.Sprintf("%s服务->%s:%d", index, addr, port))
		} else {
			thirdLiveService = append(thirdLiveService, fmt.Sprintf("%s服务->%s:%d", index, addr, port))
		}

	}
	hostname, err := os.Hostname()
	if err != nil {
		log.Println(err)
	}
	if len(thirdDownService) == 0 {
		m.thirdDownService.With(prometheus.Labels{"third_down_services": ""}).Set(0)
	} else {
		m.thirdDownService.With(prometheus.Labels{
			"third_down_services":      strings.Join(thirdDownService, "  "),
			"third_down_services_name": strings.Join(thirdDownServiceName, "  "),
			"instance":                 hostname,
		}).Set(1)
	}

	if len(thirdLiveService) == 0 {
		m.thirdLiveService.With(prometheus.Labels{"third_live_services": ""}).Set(0)
	} else {
		m.thirdLiveService.With(prometheus.Labels{"third_live_services": strings.Join(thirdLiveService, "  ")}).Set(1)
	}

	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	// http.Handle("/metrics", promHandler)
	// http.HandleFunc("/devices", getDevices)
	// http.ListenAndServe(":9000", nil)

	pMux := http.NewServeMux()
	pMux.Handle("/metrics", promHandler)

	go func() {
		log.Fatal(http.ListenAndServe(":9000", pMux))
	}()

	select {}
}
