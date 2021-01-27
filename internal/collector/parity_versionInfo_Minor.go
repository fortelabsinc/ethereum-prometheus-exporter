package collector

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

type ParityVersionInfoMinor struct {
	rpc           *rpc.Client
	desc    *prometheus.Desc
}


func NewParityVersionInfoMinor(rpc *rpc.Client) *ParityVersionInfoMinor {
	return &ParityVersionInfoMinor{
		rpc: rpc,
		desc: prometheus.NewDesc(
			"parity_versionInfo_Minor",
			"Provides information about running version of Parity minor version",
			nil,
			nil,
		),
	}
}

func (collector *ParityVersionInfoMinor) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *ParityVersionInfoMinor) Collect(ch chan<- prometheus.Metric) {
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

	value := float64(result.Version.Minor)
	ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, value)
}
