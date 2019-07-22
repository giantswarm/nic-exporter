package main

import (
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/safchain/ethtool"
)

func main() {

	// retrieve interface name fro env
	iface := os.Getenv("INTERFACE_NAME")
	if iface == "" {
		panic(errors.New("Environment variable INTERFACE_NAME can not be empty"))
	}
	_, err := ethtool.BusInfo(iface)
	if err != nil {
		panic(err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	nicCollector := newNICCollector(hostname, iface)
	prometheus.MustRegister(nicCollector)

	http.Handle("/metrics", promhttp.Handler())
	log.Info("Serving on port :10800")
	log.Fatal(http.ListenAndServe(":10800", nil))
}
