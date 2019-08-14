package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/giantswarm/exporterkit"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/giantswarm/nic-exporter/nic"
	"github.com/giantswarm/nic-exporter/nstat"
)

var (
	iface string
)

func init() {
	flag.StringVar(&iface, "iface", "eth0", "Interface name to retrieve stats from")
}

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "version" || os.Args[1] == "--help") {
		return
	}

	flag.Parse()

	var err error

	var logger micrologger.Logger
	{
		logger, err = micrologger.New(micrologger.Config{})
		if err != nil {
			panic(fmt.Sprintf("%#v\n", err))
		}
	}

	var nicCollector prometheus.Collector
	{
		c := nic.Config{
			Logger: logger,

			IFace: iface,
		}

		nicCollector, err = nic.New(c)
		if err != nil {
			panic(fmt.Sprintf("%#v\n", err))
		}
	}

	var nstatCollector prometheus.Collector
	{

		c := nstat.Config{
			Logger: logger,
		}

		nstatCollector, err = nstat.New(c)
		if err != nil {
			panic(fmt.Sprintf("%#v\n", err))
		}
	}

	var exporter *exporterkit.Exporter
	{
		c := exporterkit.Config{
			Collectors: []prometheus.Collector{
				nicCollector,
				nstatCollector,
			},
			Logger: logger,
		}

		exporter, err = exporterkit.New(c)
		if err != nil {
			panic(fmt.Sprintf("%#v\n", err))
		}
	}

	exporter.Run()
}
