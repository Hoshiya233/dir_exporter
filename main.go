package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v2"
)

func main() {
	fileSize := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "file",
			Help: "目录里各文件的大小",
		},
		[]string{
			"name",
		},
	)
	http.Handle("/metrics", promhttp.Handler())
	prometheus.MustRegister(fileSize)

	Config := getConfig()

	go collect(*fileSize, Config.Path)

	log.Fatal(http.ListenAndServe(":"+Config.Port, nil))
}

type config struct {
	//结构体里变量的名字不能和yml文件里的名字全等，这是yaml模块的坑
	Path string `yaml:"path"`
	Port string `yaml:"port"`
}

func getConfig() config {
	c := new(config)
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	if !strings.HasSuffix(c.Path, "/") {
		//如果配置文件path不是以/结尾，就加上
		c.Path = c.Path + "/"
	}

	return *c
}
