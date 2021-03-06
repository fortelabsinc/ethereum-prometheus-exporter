package collector

import (
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

type NetPeerCount struct {
	rpc  *rpc.Client
	desc *prometheus.Desc
}

func NewNetPeerCount(rpc *rpc.Client) *NetPeerCount {
	return &NetPeerCount{
		rpc: rpc,
		desc: prometheus.NewDesc(
			"net_peers",
			"the number of peers currently connected to the client",
			nil,
			nil,
		),
	}
}

func (collector *NetPeerCount) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *NetPeerCount) Collect(ch chan<- prometheus.Metric) {

	timeoutChannel := make(chan prometheus.Metric, 1)
	defer close(timeoutChannel)

	go func() {
		var result hexutil.Uint64
		start := time.Now()
		if err := collector.rpc.Call(&result, "net_peerCount"); err != nil {
			timeoutChannel <- prometheus.NewInvalidMetric(collector.desc, err)
			return
		}

		end := time.Now()
		log.Print("net_peerCount: ", end.Sub(start))

		value := float64(result)
		timeoutChannel <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, value)
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
