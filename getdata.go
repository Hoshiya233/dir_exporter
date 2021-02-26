package main

import (
	"bufio"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

func getData(dir prometheus.GaugeVec) {
	/*
		获取被监控目录的数据
	*/
	cmd := exec.Command("/bin/bash", "-c", `du -s /home/go/*`)

	//创建获取命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
	}

	//执行命令
	if err := cmd.Start(); err != nil {
		log.Println("Error:The command is err,", err)
	}

	//使用带缓冲的读取器
	outputBuf := bufio.NewReader(stdout)

	for {
		//一次获取一行,_ 获取当前行是否被读完
		output, _, err := outputBuf.ReadLine()
		if err != nil {
			// 判断是否到文件的结尾了否则出错
			if err.Error() != "EOF" {
				log.Printf("Error :%s\n", err)
			}
			break
		}
		log.Printf("%s\n", string(output))

		//将结果保存
		s := string(output)
		ss := strings.Split(s, "\t")
		key := ss[1]
		value, _ := strconv.ParseFloat(ss[0], 64)
		dir.WithLabelValues(key).Set(value)
	}

	//输出错误信息
	if err := cmd.Wait(); err != nil {
		log.Println("wait:", err.Error())
	}
}
