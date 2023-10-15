FROM golang:alpine AS builder

# 环境变量
ENV GO111MODULE=on \
GOPROXY=https://goproxy.cn,direct \
CGO_ENABLED=0 \
GOOS=linux \
GOARCH=amd64

# 移动到工作目录：/build
WORKDIR /build

# 复制项目中的 go.mod 和 go.sum文件并下载依赖信息
COPY go.mod .
COPY go.sum .
RUN go mod download

# 将代码复制到容器中
COPY . .

# 将代码编译成二进制可执行文件
RUN go build -o swapbotAA-relayer .

###################
# 创建镜像
###################
FROM debian:stretch-slim
#FROM scratch

COPY ./wait-for.sh /
COPY ./conf /conf

# 从builder镜像中把可执行文件拷贝到当前目录
COPY --from=builder /build/swapbotAA-relayer /

#RUN set -eux \
#
#    && apt-get update \
#    && apt-get install -y --no-install-recommends netcat \
#    && chmod 755 wait-for.sh

# 声明服务端口
EXPOSE 8081

# 需要运行的命令
ENTRYPOINT ["/swapbotAA-relayer", "conf/config.yaml"]
