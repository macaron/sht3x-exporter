package main

import (
	"github.com/d2r2/go-i2c"
	"github.com/d2r2/go-sht3x"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

var temperatureGauge = promauto.NewGauge(prometheus.GaugeOpts{
	Namespace: "sht3x",
	Name:      "temperature",
	Help:      "Current temperature",
})

var humidityGauge = promauto.NewGauge(prometheus.GaugeOpts{
	Namespace: "sht3x",
	Name:      "humidity",
	Help:      "Current humidity",
})

func recordMetrics() {
	go func() {
		dev, err := i2c.NewI2C(0x45, 1)
		defer dev.Close()
		if err != nil {
			log.Fatal(err)
		}
		sensor := sht3x.NewSHT3X()
		err = sensor.Reset(dev)
		if err != nil {
			log.Fatal(err)
		}
		for {
			temp, rh, err := sensor.ReadTemperatureAndRelativeHumidity(dev, sht3x.RepeatabilityHigh)
			if err != nil {
				log.Fatal(err)
			}
			temperatureGauge.Set(float64(temp))
			humidityGauge.Set(float64(rh))
			time.Sleep(5 * time.Second)
		}
	}()
}

func main() {
	recordMetrics()
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":2112", nil)
	if err != nil {
		log.Fatal(err)
	}
}
