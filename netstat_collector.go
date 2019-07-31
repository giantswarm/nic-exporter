package main

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const (
	nstat_metric_namespace = "nstat"
)

type nstatCollector struct {
	hostname     string
	nstatMetrics map[string]*prometheus.Desc
}

func newNstatCollector(hostname string) *nstatCollector {

	nstatCollector := &nstatCollector{
		hostname: hostname,
	}

	nstatOutput, err := exec.Command("/sbin/nstat", "-a", "--json").Output()
	if err != nil {
		log.Fatal(err)
	}

	var nstatObj map[string]interface{}

	err = json.Unmarshal([]byte(nstatOutput), &nstatObj)
	if err != nil {
		log.Fatal(err)
	}

	nstatMetrics, _ := nstatObj["kernel"].(map[string]interface{})

	nstatCollector.nstatMetrics = make(map[string]*prometheus.Desc)
	for label, _ := range nstatMetrics {
		fqName := prometheus.BuildFQName(nstat_metric_namespace, "", label)
		nstatCollector.nstatMetrics[label] = prometheus.NewDesc(fqName, fmt.Sprintf("Generated description for metric %#q", label), []string{"hostname"}, nil)
	}

	return nstatCollector
}

func (collector *nstatCollector) Describe(ch chan<- *prometheus.Desc) {

	for _, nstatMetric := range collector.nstatMetrics {
		ch <- nstatMetric
	}
}

func (collector *nstatCollector) Collect(ch chan<- prometheus.Metric) {

	nstatOutput, err := exec.Command("/sbin/nstat", "-a", "--json").Output()
	if err != nil {
		log.Fatal(err)
	}

	var nstatObj map[string]interface{}

	err = json.Unmarshal([]byte(nstatOutput), &nstatObj)
	if err != nil {
		log.Fatal(err)
	}

	kernelMetrics := nstatObj["kernel"].(map[string]interface{})

	for label, nstatMetric := range collector.nstatMetrics {
		ch <- prometheus.MustNewConstMetric(nstatMetric, prometheus.GaugeValue, kernelMetrics[label].(float64), collector.hostname)
	}
}
