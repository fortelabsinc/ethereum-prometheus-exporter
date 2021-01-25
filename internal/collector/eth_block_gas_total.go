package collector

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

type EthBlockGasTotal struct {
	rpc  *rpc.Client
	desc *prometheus.Desc
}

type gasResult struct {
	GasUsed hexutil.Uint64
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
	var raw json.RawMessage
	if err := collector.rpc.Call(&raw, "eth_getBlockByNumber", "latest", true ); err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}

	var result *gasResult
	if err := json.Unmarshal(raw, &result); err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}

	value:= float64(result.GasUsed)
	ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, value)
}
