FROM ubuntu:20.04 as BASE

ARG TZ=Asia/Hong_Kong
ARG DEBIAN_FRONTEND=noninteractive

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN mv /etc/apt/sources.list /etc/apt/sources.list.bak
RUN echo 'deb http://mirrors.aliyun.com/ubuntu/ focal main restricted universe multiverse\n\
deb http://mirrors.aliyun.com/ubuntu/ focal-security main restricted universe multiverse\n\
deb http://mirrors.aliyun.com/ubuntu/ focal-updates main restricted universe multiverse\n\
deb http://mirrors.aliyun.com/ubuntu/ focal-proposed main restricted universe multiverse\n\
deb http://mirrors.aliyun.com/ubuntu/ focal-backports main restricted universe multiverse'\
>> /etc/apt/sources.list

FROM BASE
RUN apt-get update && apt-get install -y \
  openssh-server sudo sed tzdata net-tools && \
  apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config && \
echo 'root:root' | chpasswd && \
useradd -rm -d /home/test -s /bin/bash -g root -G sudo -u 1000 test && \
echo 'test:test' | chpasswd

EXPOSE 22
RUN service ssh start
CMD ["/usr/sbin/sshd","-D"]
