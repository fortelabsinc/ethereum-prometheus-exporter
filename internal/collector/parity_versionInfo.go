package collector

import (
	"strconv"
	"encoding/json"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

type ParityVersionInfo struct {
	rpc           *rpc.Client
	desc    *prometheus.Desc
}

type version struct {
	Major int
	Minor int
	Patch int	
}

type versionResult struct {
	Version version
}

func NewParityVersionInfo(rpc *rpc.Client) *ParityVersionInfo {
	return &ParityVersionInfo{
		rpc: rpc,
		desc: prometheus.NewDesc(
			"parity_versionInfo",
			"Provides information about running version of Parity version",
			[]string{"version"},
			nil,
		),
	}
}

func (collector *ParityVersionInfo) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *ParityVersionInfo) Collect(ch chan<- prometheus.Metric) {
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
	var value float64 = 0
	versionValue := strconv.Itoa(result.Version.Major) + "." + strconv.Itoa(result.Version.Minor) + "." + strconv.Itoa(result.Version.Patch)
	ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, value, versionValue )
}
