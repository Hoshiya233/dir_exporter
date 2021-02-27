package main

import (
	"bufio"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func collect(fileSize prometheus.GaugeVec, path string) {
	/*
		获取被监控目录的数据
		每60秒执行一次
	*/
	shell := "du -s " + path + "*"
	log.Println(shell)
	for {
		//整体循环
		cmd := exec.Command("/bin/bash", "-c", shell)

		//创建获取命令输出管道
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
		}

		//执行命令
		log.Println("开始执行命令")
		if err := cmd.Start(); err != nil {
			log.Println("Error:The command is err,", err)
		}
		log.Println("命令执行完毕")

		//使用带缓冲的读取器
		outputBuf := bufio.NewReader(stdout)

		log.Println("开始读取结果")
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
			fileSize.WithLabelValues(key).Set(value)
		}
		log.Println("结果读取完毕")

		//输出错误信息
		if err := cmd.Wait(); err != nil {
			log.Println("wait:", err.Error())
		}

		time.Sleep(time.Second * 60)
	}
}
