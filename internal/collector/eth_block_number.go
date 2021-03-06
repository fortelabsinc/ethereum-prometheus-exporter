package collector

import (
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

type EthBlockNumber struct {
	rpc  *rpc.Client
	desc *prometheus.Desc
}

func NewEthBlockNumber(rpc *rpc.Client) *EthBlockNumber {
	return &EthBlockNumber{
		rpc: rpc,
		desc: prometheus.NewDesc(
			"eth_block_number",
			"the number of most recent block",
			nil,
			nil,
		),
	}
}

func (collector *EthBlockNumber) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *EthBlockNumber) Collect(ch chan<- prometheus.Metric) {
	timeoutChannel := make(chan prometheus.Metric, 1)
	defer close(timeoutChannel)

	go func() {
		var result hexutil.Uint64
		start := time.Now()
		if err := collector.rpc.Call(&result, "eth_blockNumber"); err != nil {
			timeoutChannel <- prometheus.NewInvalidMetric(collector.desc, err)
			return
		}

		end := time.Now()
		log.Print("eth_blockNumber: ", end.Sub(start))
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
