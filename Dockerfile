# 构建阶段
FROM golang:1.21.4-alpine AS builder

WORKDIR /build

# 安装构建依赖
RUN apk add --no-cache \
    git \
    gcc \
    musl-dev \
    make

# 复制源代码和依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o isscan main.go

# 运行阶段
FROM alpine:3.19

WORKDIR /app

# 安装运行时依赖
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl \
    wget \
    bind-tools

# 从构建阶段复制二进制文件
COPY --from=builder /build/isscan /usr/local/bin/

# 创建非 root 用户
RUN adduser -D -g '' isscan && \
    mkdir -p /app/logs && \
    chown -R isscan:isscan /app

USER isscan

# 暴露端口
EXPOSE 8888

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget -q --spider http://localhost:8888/api/ping || exit 1

# 启动应用
ENTRYPOINT ["/usr/local/bin/isscan"]
CMD ["--config", "/app/config.yaml"] 