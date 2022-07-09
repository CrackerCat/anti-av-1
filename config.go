package main

type config struct {
	//生成shellcode
	loader    string //构造shellcode的方式
	shellcode string //bin格式shellcode文件
	os        string //windows,linux
	sign      string //签名信息

	//加密 shellcode
	crypt bool
}
