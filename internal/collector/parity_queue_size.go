package collector

import (
	"log"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

type ParityQueueSize struct {
	rpc  *rpc.Client
	desc *prometheus.Desc
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
	var result *[]interface{}
	start := time.Now()
	if err := collector.rpc.Call(&result, "parity_allTransactions"); err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}
	end := time.Now()

	log.Print("parity_allTransactions: ", end.Sub(start))
	value := float64(len(*result))

	ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, value)
}
