package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v2"
)

//VERSION 版本号
var VERSION = "0.1"
var isPrintVersion = false

type config struct {
	//结构体里变量的名字不能和yml文件里的名字全等，这是yaml模块的坑
	Path string `yaml:"path"`
	Port string `yaml:"port"`
}

func main() {
	Config := getConfig()
	//先读取配置文件里的参数，再获取命令行参数，因此命令行配置优先级更高
	flag.StringVar(&Config.Path, "path", Config.Path, "被监控的目录")
	flag.StringVar(&Config.Port, "port", Config.Port, "指定开放的端口号")
	flag.BoolVar(&isPrintVersion, "v", false, "显示版本号，然后退出")
	flag.Parse()

	if isPrintVersion {
		fmt.Println("version:", VERSION)
		os.Exit(0)
	}
	if !strings.HasSuffix(Config.Path, "/") {
		//如果配置文件path不是以/结尾，就加上
		Config.Path = Config.Path + "/"
	}

	log.Println("path:", Config.Path)
	log.Println("port:", Config.Port)

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

	go collect(*fileSize, Config.Path)

	log.Fatal(http.ListenAndServe(":"+Config.Port, nil))
}

func getConfig() config {
	c := new(config)
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		c.Path = "/var/log/"
		c.Port = "8816"
	}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return *c
}
