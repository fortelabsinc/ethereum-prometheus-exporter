package collector

import (
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

type EthEarliestBlockTransactions struct {
	rpc  *rpc.Client
	desc *prometheus.Desc
}

func NewEthEarliestBlockTransactions(rpc *rpc.Client) *EthEarliestBlockTransactions {
	return &EthEarliestBlockTransactions{
		rpc: rpc,
		desc: prometheus.NewDesc(
			"eth_earliest_block_transactions",
			"the number of transactions in an earliest block",
			nil,
			nil,
		),
	}
}

func (collector *EthEarliestBlockTransactions) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *EthEarliestBlockTransactions) Collect(ch chan<- prometheus.Metric) {
	var result hexutil.Uint64
	start := time.Now()
	if err := collector.rpc.Call(&result, "eth_getBlockTransactionCountByNumber", "earliest"); err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}
	end := time.Now()
	log.Print("eth_getBlockTransactionCountByNumber earliest: ", end.Sub(start))

	value := float64(result)
	ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, value)
}
