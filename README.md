# vpsgo

极简 VPS 服务远程管理工具。

## 安装


## 使用

要求 VPS 安装 Linux 系统。当前支持 Ubuntu 20.04。

### 注册 VPS

注册 VPS 服务，同时配置 ssh 无密码证书登录。

使用条件: 有系统 root 密码，且允许 ssh 密码登录。

```sh
vps reg VPS_NAME IP_ADDR [PORT]
```

**参数说明：**
- VPS_NAME - 自定义 VPS 名称，必须唯一。
- IP_ADDR - VPS 服务器 IP 地址。
- PORT - ssh 服务端口号，默认 22。

注：命令行中大写的部分表示参数，带方括号的是可选参数，下同。

### 列表 VPS

显示当前已经注册的 VPS 列表。

```sh
vps list
```

### 设置默认 VPS

设置默认操作的 VPS 服务器。

```sh
vps use VPS_NAME
```

### 安装 VPS 服务

在默认 VPS 上安装服务软件。

```sh
vps install nginx
```

### 查看 VPS 状态

查看指定或默认 VPS 服务当前状态。

```sh
vps status [VPS_NAME]
```
