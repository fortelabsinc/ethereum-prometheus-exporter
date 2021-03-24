package collector

import (
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

type EthGasPrice struct {
	rpc  *rpc.Client
	desc *prometheus.Desc
}

func NewEthGasPrice(rpc *rpc.Client) *EthGasPrice {
	return &EthGasPrice{
		rpc: rpc,
		desc: prometheus.NewDesc(
			"eth_gas_price",
			"the current price per gas in wei",
			nil,
			nil,
		),
	}
}

func (collector *EthGasPrice) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *EthGasPrice) Collect(ch chan<- prometheus.Metric) {
	var result hexutil.Big
	start := time.Now()
	time.AfterFunc(3*time.Second, func() {
		if err := collector.rpc.Call(&result, "eth_gasPrice"); err != nil {
			errorEnd := time.Now()
			log.Print("error eth_gasPrice: ", errorEnd.Sub(start))
			ch <- prometheus.NewInvalidMetric(collector.desc, err)
			return
		}
	})
	end := time.Now()
	log.Print("eth_gasPrice: ", end.Sub(start))
	i := (*big.Int)(&result)
	value, _ := new(big.Float).SetInt(i).Float64()
	ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, value)
}
