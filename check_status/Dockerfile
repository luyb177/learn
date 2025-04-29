# 阶段 1：构建 Go 应用（Builder）
FROM golang:1.23.4 AS builder

WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct

# 1. 复制依赖文件（利用 Docker 缓存优化）
COPY go.mod go.sum ./
RUN go mod download

# 2. 复制所有源码
COPY . .

# 3. 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/check_status .

# 阶段 2：运行环境（最小化镜像）
FROM alpine:latest

WORKDIR /app

# ✅ 安装 tzdata，用于时区支持
RUN apk add --no-cache tzdata

# 复制可执行文件和配置
COPY --from=builder /app/check_status /app/check_status
COPY config/config.yml /app/config/config.yml

# 设置时区为上海（可选，但更明确）
ENV TZ=Asia/Shanghai

EXPOSE 7777

CMD ["/app/check_status"]
