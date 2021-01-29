package collector

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

type ParityQueueSize struct {
	rpc           *rpc.Client
	desc    *prometheus.Desc
}


func NewParityQueueSize(rpc *rpc.Client) *ParityQueueSize {
	return &ParityQueueSize{
		rpc: rpc,
		desc: prometheus.NewDesc(
			"parity_queue_size",
			"Retrun size of the queue",
			nil,
			nil,
		),
	}
}

func (collector *ParityQueueSize) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *ParityQueueSize) Collect(ch chan<- prometheus.Metric) {
	var raw json.RawMessage
	if err := collector.rpc.Call(&raw, "parity_accountsInfo"); err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc , err)
		return
	}

	var result *[]interface{}
	if err := json.Unmarshal(raw, &result); err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}

	value := float64(len(*result))

	ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, value)
}
