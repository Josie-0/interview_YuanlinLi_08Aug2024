# 使用官方 Go 镜像作为构建环境
FROM golang:1.22 AS builder

# 设置工作目录
WORKDIR /app

# 复制 go mod 文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

# 暴露应用的端口
EXPOSE 8080

CMD ["/docker-gs-ping"]