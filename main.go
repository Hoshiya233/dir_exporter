package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	dir := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dir",
			Help: "目录里各文件的大小",
		},
		[]string{
			"name",
		},
	)
	http.Handle("/metrics", promhttp.Handler())
	prometheus.MustRegister(dir)

	go update(*dir)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func update(dir prometheus.GaugeVec) {
	//每60秒调用一次getData函数，获取数据
	for {
		getData(dir)
		time.Sleep(time.Second * 60)
	}
}
