package main

import (
	"bytes"
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ENV vars config
type settings struct {
	WebURL string `envconfig:"WEB_URL" default:"localhost:8080"`
}

var webJSON = []byte(`{
		"object": {
				"param": "value"
		},
		"param": "value"
}`)

func main() {
	/*
		SETUP:
			1) Load settings from ENV_VARS
			2) Run Prometheus in goroutine
	*/
	log.Println("Starting...")
	var s settings
	err := envconfig.Process("envs", &s)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Loaded settings:")
	log.Println(s)

	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "WebStatus",
		Help: "Status of specified Web API.  0=Healthy; 1=Failed;",
	})
	prometheus.MustRegister(gauge)
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
	log.Println("Started Prometheus...")
	log.Println("Started Web-Health-Check Service")

	/*
		MAIN LOOP:
			Tests web api endpoint.
			If fails, sets gauge to 1 and retries in 1 minute
			If passes, sleep for 5 minutes and try again
	*/
	var failed bool
	for {
		failed = testWebAPI(s)

		if failed {
			gauge.Set(1)
			time.Sleep(1 * time.Minute)
		} else {
			gauge.Set(0)
			time.Sleep(5 * time.Minute)
		}
	}
}

func testWebAPI(s settings) (failed bool) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Post(s.WebURL, "application/json", bytes.NewBuffer(webJSON))
	if err != nil {
		log.Println("err: " + err.Error())
		return true
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return true
	}
	return false
}
