package collector

import (
	 "log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

type EthBlockGasTotal struct {
	rpc  *rpc.Client
	desc *prometheus.Desc
}

type gasResult struct {
	GasLimit hexutil.Uint64
}

func NewEthBlockGasTotal(rpc *rpc.Client) *EthBlockGasTotal {
	return &EthBlockGasTotal{
		rpc: rpc,
		desc: prometheus.NewDesc(
			"eth_block_gas_total",
			"The total gas used in the block, given by the field gasUsed",
			nil,
			nil,
		),
	}
}

func (collector *EthBlockGasTotal) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *EthBlockGasTotal) Collect(ch chan<- prometheus.Metric) {
	var result *gasResult
	if err := collector.rpc.Call(&result, "eth_getBlockByNumber", "latest", true ); err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}
	log.Println(result)
	value := float64(result.GasLimit)
	ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, value)
}
