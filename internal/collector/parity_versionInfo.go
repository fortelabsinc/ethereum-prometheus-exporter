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
			[]string{"version"},
			nil,
		),
	}
}

func (collector *ParityVersionInfo) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *ParityVersionInfo) Collect(ch chan<- prometheus.Metric) {
	timeoutChannel := make(chan prometheus.Metric, 1)
	defer close(timeoutChannel)

	go func() {
		var result string
		start := time.Now()
		if err := collector.rpc.Call(&result, "web3_clientVersion"); err != nil {
			errorEnd := time.Now()
			log.Print("error web3_clientVersion: ", errorEnd.Sub(start))
			timeoutChannel <- prometheus.NewInvalidMetric(collector.desc, err)
			return
		}
		end := time.Now()

		log.Print("web3_clientVersion: ", end.Sub(start))
		value := float64(1)
		timeoutChannel <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, value, result)
	}()
	timer := time.NewTimer(2 * time.Second)
	defer timer.Stop()

	select {
	case rpcCallResult := <-timeoutChannel:
		ch <- rpcCallResult
	case <-timer.C:
		log.Print("net_peerCount Timed out")
		ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, 0)
	}

}
