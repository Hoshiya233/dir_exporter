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
var VERSION = "0.1.2"
var isPrintVersion = false

//ConfigFilePath 配置文件路径
var ConfigFilePath = "./config.yml"

type config struct {
	//结构体里变量的名字不能和yml文件里的名字全等，这是yaml模块的坑
	Paths []string `yaml:"paths"`
	Port  string   `yaml:"port"`
}

func init() {
	// flag.StringVar(&Config.Port, "port", Config.Port, "指定开放的端口号")
	// flag.StringVar(&Config.Path, "path", Config.Path, "被监控的目录")
	flag.StringVar(&ConfigFilePath, "c", ConfigFilePath, "配置文件路径，绝对路径")
	flag.BoolVar(&isPrintVersion, "v", false, "显示版本号，然后退出")
	flag.Parse()
	if isPrintVersion {
		fmt.Println("version:", VERSION)
		os.Exit(0)
	}
}

func main() {
	Config := getConfig()
	for i, Path := range Config.Paths {
		if !strings.HasSuffix(Path, "/") {
			//如果配置文件path不是以/结尾，就加上
			Config.Paths[i] = Path + "/"
		}
	}

	log.Println("port:", Config.Port)
	log.Println("path:", Config.Paths)

	fileSize := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "filesize",
			Help: "目录里各文件的大小",
		},
		[]string{
			"name",
		},
	)
	http.Handle("/metrics", promhttp.Handler())
	prometheus.MustRegister(fileSize)

	go collect(*fileSize, Config.Paths)

	log.Fatal(http.ListenAndServe(":"+Config.Port, nil))
}

func getConfig() config {
	c := new(config)
	yamlFile, err := ioutil.ReadFile(ConfigFilePath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		c.Port = "8816"
		c.Paths = []string{"/var/log/"}
	}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return *c
}
