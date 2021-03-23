package collector

import (
	"log"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

type ParityVersionInfo struct {
	rpc  *rpc.Client
	desc *prometheus.Desc
}

func NewParityVersionInfo(rpc *rpc.Client) *ParityVersionInfo {
	return &ParityVersionInfo{
		rpc: rpc,
		desc: prometheus.NewDesc(
			"parity_versionInfo",
			"Provides information about running version of Parity version",
			nil,
			nil,
		),
	}
}

func (collector *ParityVersionInfo) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *ParityVersionInfo) Collect(ch chan<- prometheus.Metric) {
	var result string
	start := time.Now()
	if err := collector.rpc.Call(&result, "web3_clientVersion"); err != nil {
		errorEnd := time.Now()
		log.Print("error web3_clientVersion: ", errorEnd.Sub(start))
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}
	end := time.Now()

	log.Print("web3_clientVersion: ", end.Sub(start))
	log.Print(result)
	// if err := json.Unmarshal(raw, &result); err != nil {
	// 	ch <- prometheus.NewInvalidMetric(collector.desc, err)
	// 	return
	// }
	value := float64(1)
	// versionValue := strconv.Itoa(result.Version.Major) + "." + strconv.Itoa(result.Version.Minor) + "." + strconv.Itoa(result.Version.Patch)
	ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, value)
}
