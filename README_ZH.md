# vpsgo - 极简 VPS 服务远程管理工具。

[English Version](README.md)

## 愿景

vpsgo 想要做最简单的 VPS 服务器管理工具。

## 简介

众所周知，VPS 是一个远程服务器。通常你需要使用 ssh 远程登录来对它进行配置和管理。你必须要熟悉各种操作命令，了解服务器系统和相关特性，具备脚本编程技能才能胜任服务器管理工作。而且，相关工作十分繁琐，费时费力而又充满了风险。现在，所有这些你都不需要关心，因为有了 vpsgo。

vpsgo 整合服务器管理的最佳实践，使用云服务开发语言 Golang 开发。vpsgo 致力于带给你一个轻便、简约、高效的服务器管理体验，减轻你的心智负担，释放你的生产力。

## 安装

### 通过源码安装

#### 安装 Go

vpsgo 需要 Go 1.15 编译, 请参考 [官方文档](https://golang.org/doc/install) 安装。

#### 安装 Docker

可以使用容器练习 vpsgo 的用法。请参考 [官方文档](https://docs.docker.com/engine/install/) 安装 Docker。

#### 安装构建工具

**Linux**

```sh
# Ubuntu or Debian
sudo apt-get install build-essential git

# Fedora or CentOS
sudo yum groupinstall "Development Tools"
# or
sudo yum install gcc git -y
```

**macOS**

打开 "终端"（位于 应用/工具）

在 终端 窗口执行命令 `xcode-select --install`

在弹出的窗口中点击 Install 安装，选择 同意 Terms of Service 服务条款。

**Windows**

推荐使用 [WSL 2 (安装Ubuntu)](https://docs.microsoft.com/en-us/windows/wsl/install-win10) 和 [Windows Terminal](https://docs.microsoft.com/en-us/windows/terminal/get-started)。否则，请安装 [Git for Windows](https://gitforwindows.org/) 以获得 `Git Bash` 和许多 Linux 命令. 然后安装 [Chocolatecy](https://chocolatey.org/install) 或 [Scoop](https://scoop.sh/) 以便安装相关扩展工具。

```sh
choco install make
# or
scoop install make
```

#### 编译 vpsgo

```sh
# Clone the repository to the "vpsgo" subdirectory
git clone --depth 1 https://github.com/appmind/vpsgo.git vpsgo

# Change working directory
cd vpsgo

# Compile the main program, dependencies will be downloaded at this step
make install
```

#### 测试安装

```sh
# Start the container server
make docker-up

# Check the container server
vps ping 127.0.0.1 -p 22 -u root -P root

# View usage help
vps help
```

## 协议

本项目采用 Apache License 2.0 开源授权协议。完整的许可证文本，请参见 [LICENSE](LICENSE) 文件。
