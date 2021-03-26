package collector

import (
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

type EthHashrate struct {
	rpc  *rpc.Client
	desc *prometheus.Desc
}

func NewEthHashrate(rpc *rpc.Client) *EthHashrate {
	return &EthHashrate{
		rpc: rpc,
		desc: prometheus.NewDesc(
			"eth_hashrate",
			"the number of hashes per second that the node is mining with",
			nil,
			nil,
		),
	}
}

func (collector *EthHashrate) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *EthHashrate) Collect(ch chan<- prometheus.Metric) {
	timeoutChannel := make(chan prometheus.Metric, 1)
	defer close(timeoutChannel)

	go func() {

		var result hexutil.Uint64
		start := time.Now()
		if err := collector.rpc.Call(&result, "eth_hashrate"); err != nil {
			timeoutChannel <- prometheus.NewInvalidMetric(collector.desc, err)
			return
		}
		end := time.Now()
		log.Print("eth_hashrate: ", end.Sub(start))
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
