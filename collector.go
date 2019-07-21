package main

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/safchain/ethtool"
)

const (
	nic_metric_namespace = "nic"
)

type nicCollector struct {
	nicMetrics map[string]*prometheus.Desc
	iface      string
}

func newNICCollector(iface string) *nicCollector {
	nicCollector := &nicCollector{
		iface: iface,
	}

	nicStats, err := ethtool.Stats(iface)
	if err != nil {
		panic(err)
	}
	nicCollector.nicMetrics = make(map[string]*prometheus.Desc)
	for label, _ := range nicStats {
		fqName := prometheus.BuildFQName(nic_metric_namespace, "", label)
		nicCollector.nicMetrics[label] = prometheus.NewDesc(fqName, fmt.Sprintf("Generated description for metric %#q", label), []string{"iface"}, nil)
	}

	return nicCollector
}

func (collector *nicCollector) Describe(ch chan<- *prometheus.Desc) {

	for _, nicMetric := range collector.nicMetrics {
		ch <- nicMetric
	}
}

func (collector *nicCollector) Collect(ch chan<- prometheus.Metric) {

	nicStats, _ := ethtool.Stats(collector.iface)

	for label, nicMetric := range collector.nicMetrics {
		ch <- prometheus.MustNewConstMetric(nicMetric, prometheus.GaugeValue, float64(nicStats[label]), collector.iface)
	}
}
