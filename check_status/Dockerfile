# ---------------------
# 阶段 1：构建 Go 应用
# ---------------------
FROM golang:1.23.4 AS builder

WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct

# 1. 复制依赖文件并缓存依赖
COPY go.mod go.sum ./
RUN go mod download

# 2. 复制源代码
COPY . .

# 3. 构建可执行文件（使用 -tags timetzdata 内嵌时区数据更稳妥）
RUN CGO_ENABLED=0 GOOS=linux go build -tags timetzdata -ldflags="-w -s" -o /app/check_status .

# ---------------------
# 阶段 2：运行环境
# ---------------------
FROM alpine:latest

WORKDIR /app

# ✅ 安装 tzdata（提供系统时区信息，虽然使用 timetzdata tag 已足够稳妥）
RUN apk add --no-cache tzdata

# ✅ 拷贝 zoneinfo 数据（Go runtime 有些版本仍需）
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# ✅ 设置系统环境时区
ENV TZ=Asia/Shanghai

# ✅ 拷贝可执行文件和配置文件
COPY --from=builder /app/check_status /app/check_status
COPY config/config.yml /app/config/config.yml

# ✅ 暴露应用端口
EXPOSE 7777

# ✅ 启动应用
CMD ["/app/check_status"]
