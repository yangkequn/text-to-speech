# 使用 golang:buster 镜像作为构建环境
FROM golang:1.20-buster as builder

# 安装编译依赖
RUN apt-get update && \
    apt-get install -y libasound2 wget build-essential 

# 下载并编译指定版本的 OpenSSL
RUN wget -O - https://www.openssl.org/source/openssl-1.1.1u.tar.gz | tar zxf - && \
    cd openssl-1.1.1u && \
    ./config --prefix=/usr/local && \
    make -j $(nproc) && \
    make install_sw install_ssldirs && \
    cd .. && rm -rf openssl-1.1.1u*

WORKDIR /app
COPY . .

# 设置编译环境变量
ENV CGO_ENABLED=1
ENV SPEECHSDK_ROOT="/app/speechsdk"
ENV ARCHITECTURE="x64"
ENV CGO_CFLAGS="-I${SPEECHSDK_ROOT}/include/c_api"
ENV CGO_LDFLAGS="-L${SPEECHSDK_ROOT}/lib/${ARCHITECTURE} -lMicrosoft.CognitiveServices.Speech.core"
ENV LD_LIBRARY_PATH="${SPEECHSDK_ROOT}/lib/${ARCHITECTURE}:${LD_LIBRARY_PATH}"

# 构建应用
RUN go build -o tts .

# 使用 Ubuntu 20.04 LTS 作为运行环境
FROM ubuntu:20.04

# 避免在docker build过程中出现与tzdata有关的交互式提示
ENV DEBIAN_FRONTEND=noninteractive

# 安装运行时依赖
RUN apt-get update && \
    apt-get install -y libasound2 ca-certificates tzdata && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# 从构建镜像中拷贝编译好的 OpenSSL 和其他必要的文件
COPY --from=builder /usr/local/ssl /usr/local/ssl
COPY --from=builder /usr/local/bin /usr/local/bin
COPY --from=builder /usr/local/lib /usr/local/lib
COPY --from=builder /app/tts /app/tts
COPY --from=builder /app/speechsdk/lib/x64/*.so /usr/local/lib/

# 设置环境变量
ENV SSL_CERT_DIR=/etc/ssl/certs
ENV SSL_CERT_FILE=/etc/ssl/certs/ca-certificates.crt
ENV LD_LIBRARY_PATH="/usr/local/lib:${LD_LIBRARY_PATH}"

# 更新共享库缓存
RUN ldconfig

# 运行应用
CMD ["/app/tts"]