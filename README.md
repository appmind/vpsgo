# vpsgo

`vpsgo` is a minimalist VPS services management tool.

[简体中文](README_ZH.md)

## 🔮 Vision

`vpsgo` wants to be the simplest VPS service management tool. To manage VPS and servers, you can use various commands, web-based control panels, or use vpsgo.

## 📡 Overview

`vpsgo` is a CLI tool developed in golang that helps you manage VPS services more simply and easily.

As we all know, VPS is a remote server. Usually you need to use ssh remote login to configure and manage it. You must be familiar with various operating commands, understand the server system and related characteristics, and have script programming skills to be competent in server management. Moreover, the related work is very tedious, time-consuming, laborious and full of risks. Now, you don’t need to care about all of these because of vpsgo.

vpsgo integrates the best practices of server management and is developed using the cloud service development language Golang. vpsgo is committed to bringing you a light, simple, and efficient server management experience, reducing your mental burden and releasing your productivity.

## 📜 Installation

### Install from source

#### Installing Go

vpsgo requires Go 1.15 to compile, please refer to the [official documentation](https://golang.org/doc/install) for how to install Go in your system.

#### Installing Docker

You can use containers to practice `vpsgo` usage. Please refer to the [official documentation](https://docs.docker.com/engine/install/) for how to install Docker in your system.

#### Installing build tools

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

Open "Terminal" (it is located in Applications/Utilities)

In the terminal window, run the command `xcode-select --install`

In the windows that pops up, click Install, and agree to the Terms of Service.

**Windows**

It is recommended to use [WSL 2 (install Ubuntu)](https://docs.microsoft.com/en-us/windows/wsl/install-win10) and [Windows Terminal](https://docs.microsoft.com/en-us/windows/terminal/get-started). Otherwise, please install [Git for Windows](https://gitforwindows.org/) to get `Git Bash` and many Linux commands. Then install [Chocolatecy](https://chocolatey.org/install) or [Scoop](https://scoop.sh/) to install some extended commands.

```sh
choco install make
# or
scoop install make
```

#### Compile vpsgo

```sh
# Clone the repository to the "vpsgo" subdirectory
git clone --depth 1 https://github.com/appmind/vpsgo.git vpsgo

# Change working directory
cd vpsgo

# Compile the main program, dependencies will be downloaded at this step
make install
```

#### Test Installation

```sh
# Start the container server
make docker-up

# Check the container server
vps ping 127.0.0.1 -p 22 -u root -P root

# View usage help
vps help
```

## ⚖️ License

This project is under the Apache License 2.0. See the [LICENSE](LICENSE) file for the full license text.
