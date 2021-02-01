package collector

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

type ParityBlockGapStatus struct {
	rpc           *rpc.Client
	desc    *prometheus.Desc
}

type blockGapStatusResult struct {
	BlockGap []hexutil.Uint64
}

func NewParityBlockGapStatus(rpc *rpc.Client) *ParityBlockGapStatus {
	return &ParityBlockGapStatus{
		rpc: rpc,
		desc: prometheus.NewDesc(
			"parity_block_gap_status",
			"Retrun chain block gap status",
			nil,
			nil,
		),
	}
}

func (collector *ParityBlockGapStatus) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *ParityBlockGapStatus) Collect(ch chan<- prometheus.Metric) {
	var raw json.RawMessage
	if err := collector.rpc.Call(&raw, "parity_chainStatus"); err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc , err)
		return
	}

	var result *blockGapStatusResult
	if err := json.Unmarshal(raw, &result); err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}

	value := float64(len(result.BlockGap))
	if value > 0 {
		value = 0
	}else {
		value = 1
	}
	ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, value)
}
