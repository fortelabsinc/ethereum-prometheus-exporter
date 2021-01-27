package collector

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

type ParityVersionInfoPatch struct {
	rpc           *rpc.Client
	desc    *prometheus.Desc
}


func NewParityVersionInfoPatch(rpc *rpc.Client) *ParityVersionInfoPatch {
	return &ParityVersionInfoPatch{
		rpc: rpc,
		desc: prometheus.NewDesc(
			"parity_versionInfo_Patch",
			"Provides information about running version of Parity Patch version",
			nil,
			nil,
		),
	}
}

func (collector *ParityVersionInfoPatch) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *ParityVersionInfoPatch) Collect(ch chan<- prometheus.Metric) {
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

	value := float64(result.Version.Patch)
	ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, value)
}
