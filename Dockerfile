FROM ubuntu:20.04

ARG TZ=Asia/Hong_Kong
ARG DEBIAN_FRONTEND=noninteractive

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN apt-get update && apt-get install -y openssh-server sudo sed tzdata && \
apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config && \
echo 'root:root' | chpasswd && \
useradd -rm -d /home/vpsgo -s /bin/bash -g root -G sudo -u 1000 test && \
echo 'test:test' | chpasswd
RUN service ssh start

EXPOSE 22
CMD ["/usr/sbin/sshd","-D"]
