# Auto-generated Dockerfile - DO NOT EDIT

FROM alpine:latest

# ==== gcc.lang ====
# Install gcc
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
    apk update

# 安装 GCC 7.4.0、8.3.0、9.2.0
RUN apk add --no-cache gcc=7.4.0-r0 gcc=8.3.0-r0 gcc=9.2.0-r0

# 安装 GCC 7.4.0、8.3.0、9.2.0
RUN apk add --no-cache \
    gcc7=7.4.0-r1 \
    g++7=7.4.0-r1 \
    gcc8=8.3.0-r0 \
    g++8=8.3.0-r0 \
    gcc9=9.2.0-r0 \
    g++9=9.2.0-r0

RUN ls -l /usr/bin/gcc* /usr/bin/g++*

# ==== jvm.lang ====
RUN mkdir -p /tmp/java && mkdir -p /usr/lib/jvm

# Install jdk21.0.2

RUN wget https://download.java.net/java/GA/jdk21.0.2/f2283984656d49d69e91c558476027ac/13/GPL/openjdk-21.0.2_linux-x64_bin.tar.gz  -P /tmp/java
RUN tar xfvz /tmp/java/openjdk-21.0.2_linux-x64_bin.tar.gz --directory /usr/lib/jvm
RUN ls /usr/lib/jvm/jdk-21.0.2 # To confirm that the directory has the jdk-21.0.2 file!

RUN rm /tmp/java/*

# Install jdk8

RUN wget https://repo.huaweicloud.com/java/jdk/8u202-b08/jdk-8u202-linux-x64.tar.gz -P /tmp/java
RUN tar xfvz /tmp/java/jdk-8u202-linux-x64.tar.gz --directory /usr/lib/jvm
RUN ls /usr/lib/jvm/jdk1.8.0_202 # To confirm that the directory has the jdk1.8.0_202 file!

RUN rm /tmp/java/*

# ==== python.lang ====
RUN mkdir -p /usr/local/python3.8 /usr/local/python2.7

# ----------------------------
# 安装 Python 3.8.1
# ----------------------------
# 下载 Python 3.8.1 源码
RUN wget https://www.python.org/ftp/python/3.8.1/Python-3.8.1.tar.xz -P /tmp

# 解压并编译安装
RUN tar -xf /tmp/Python-3.8.1.tar.xz -C /tmp
RUN cd /tmp/Python-3.8.1 && \
    ./configure \
    --prefix=/usr/local/python3.8 \
    --enable-optimizations \
    --with-ssl-default-suites=openssl && \
    make -j$(nproc) && \
    make install

# ----------------------------
# 安装 Python 2.7.17
# ----------------------------
# 下载 Python 2.7.17 源码
RUN wget https://www.python.org/ftp/python/2.7.17/Python-2.7.17.tar.xz -P /tmp

# 解压并编译安装（需修复 musl libc 兼容性）
RUN tar -xf /tmp/Python-2.7.17.tar.xz -C /tmp
# 应用补丁修复 Alpine/musl 兼容性问题
RUN sed -i 's/\(#ifdef HAVE_GETENTROPY\)/\1\n#undef HAVE_GETENTROPY/' /tmp/Python-2.7.17/Modules/_randommodule.c
RUN cd /tmp/Python-2.7.17 && \
    ./configure \
    --prefix=/usr/local/python2.7 \
    --enable-unicode=ucs4 \
    --with-ensurepip=install && \
    make -j$(nproc) && \
    make install

# ----------------------------
# 清理临时文件
# ----------------------------
RUN rm -rf /tmp/*

# 验证安装
RUN /usr/local/python3.8/bin/python3 --version && \
    /usr/local/python2.7/bin/python --version

