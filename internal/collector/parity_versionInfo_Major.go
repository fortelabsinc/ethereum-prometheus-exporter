package collector

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

type ParityVersionInfoMajor struct {
	rpc           *rpc.Client
	desc    *prometheus.Desc
}

type version struct {
	Major uint
	Minor uint
	Patch uint	
}

type versionResult struct {
	Version version
}

func NewParityVersionInfoMajor(rpc *rpc.Client) *ParityVersionInfoMajor {
	return &ParityVersionInfoMajor{
		rpc: rpc,
		desc: prometheus.NewDesc(
			"parity_versionInfo_Major",
			"Provides information about running version of Parity major version",
			nil,
			nil,
		),
	}
}

func (collector *ParityVersionInfoMajor) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *ParityVersionInfoMajor) Collect(ch chan<- prometheus.Metric) {
	var raw json.RawMessage
	if err := collector.rpc.Call(&raw, "parity_versionInfo"); err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc , err)
		return
	}

	var result *versionResult
	if err := json.Unmarshal(raw, &result); err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}

	value := float64(result.Version.Major)
	ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, value)
}
