---
layout: post
title: TFirewall
category: 工具
tags: firewall,Golang
keywords: firewall,Golang
---


TFirewall是测试已控主机哪些端口可以出网的工具.工具主体使用golang开发,工具包含客户端及服务端,适用于windows,linux,x64,x86.
# 使用方法

* 将```server_linux_x64```上传到VPS服务器(假设VPS的IP地址10.10.10.10)
* 直接运行```.\server_linux_x64```启动监听,工具会自动监听常用的TOP10端口,或通过 ```.\server_linux_x64 20-23,53,80```自定义监听端口
* 服务端打印如下信息表示监听成功
```
root@taZ:~# ./server_linux_x64 20-23,53,80
Server listening:  [20 21 22 23 53 80]
```
* 通过WEBShell管理工具将```client_win_x86.exe```上传到已控服务器
* 运行```client_win_x86.exe 10.10.10.10```测试TOP10端口哪些可以出网(服务端对应命令```.\server_linux_x64```)
* 运行```lient_win_x86.exe 10.10.10.10 20-23,53,80```探测指定端口(服务端对应命令```.\server_linux_x64 20-23,53,80```)
* 观察服务端输出即可查看客户端可以连接服务端哪些端口,使用什么协议

```
root@iZj6cbux9hc5eo9oyud2taZ:~# ./server_linux_x64 23,80-82
Server listening:  [23 80 81]
RecvTCP On 172.17.20.209:23 From 175.132.138.137:16918 
RecvTCP On 172.17.20.209:80 From 175.132.138.137:50641 
RecvTCP On 172.17.20.209:81 From 175.132.138.137:38747 
RecvUDP On [::]:23 From 175.167.138.137:28822
RecvUDP On [::]:80 From 175.167.138.137:26831
RecvUDP On [::]:81 From 175.167.138.137:26832
```


# 已测试
## server
* ubuntu 18
## client
* Windows 10 
* ubuntu 16
* kali




