# 使用官方 Go 运行时作为基础镜像
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的包
RUN apk add --no-cache git ca-certificates tzdata

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 使用轻量级的alpine镜像作为运行环境
FROM alpine:latest

# 安装ca证书和时区数据
RUN apk --no-cache add ca-certificates tzdata curl

# 设置工作目录
WORKDIR /root/

# 从builder阶段复制二进制文件
COPY --from=builder /app/main .

# 创建必要的目录
RUN mkdir -p uploads apikey logs

# 复制配置文件
COPY --from=builder /app/apikey ./apikey/

# 暴露端口
EXPOSE 8080

# 设置环境变量
ENV GIN_MODE=release

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# 运行应用
CMD ["./main"]
