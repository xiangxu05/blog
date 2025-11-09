# 第一阶段：构建阶段
FROM golang:1.23-alpine AS builder

# 配置 Alpine 镜像源（使用国内镜像加速，优先使用较快的镜像）
RUN sed -i 's/dl-cdn.alpinelinux.org/mirror.nju.edu.cn/g' /etc/apk/repositories || \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 设置工作目录
WORKDIR /build

# 配置 Go 代理（使用国内镜像加速）
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn

# 安装必要的构建工具（最小化依赖，尝试不使用 CGO）
# 如果使用 modernc.org/sqlite（纯 Go 实现），则不需要 gcc 和 musl-dev
RUN apk add --no-cache git ca-certificates tzdata

# 复制 go.mod 和 go.sum（如果存在）
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
# 尝试禁用 CGO（modernc.org/sqlite 是纯 Go 实现，不需要 CGO）
# 如果构建失败，可以改回 CGO_ENABLED=1 并安装 gcc musl-dev sqlite-dev
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o blog ./main.go

# 第二阶段：运行阶段
FROM alpine:latest

# 配置 Alpine 镜像源（使用国内镜像加速，优先使用较快的镜像）
RUN sed -i 's/dl-cdn.alpinelinux.org/mirror.nju.edu.cn/g' /etc/apk/repositories || \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 安装必要的运行时依赖
# sqlite 运行时需要 sqlite-libs
RUN apk add --no-cache ca-certificates tzdata sqlite-libs

# 设置时区
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 创建非 root 用户
RUN addgroup -g 1000 bloguser && \
    adduser -D -u 1000 -G bloguser bloguser

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /build/blog .

# 复制 web 静态文件目录
COPY --from=builder /build/web ./web

# 复制启动脚本
COPY scripts/docker-entrypoint.sh /app/
RUN chmod +x /app/docker-entrypoint.sh && \
    chown bloguser:bloguser /app/docker-entrypoint.sh

# 声明 data 目录为挂载点（用于持久化数据库文件）
VOLUME ["/app/data"]

# 切换到非 root 用户
USER bloguser

# 暴露端口
EXPOSE 80

# 设置环境变量
ENV GIN_MODE=release

# 启动应用
CMD ["/app/docker-entrypoint.sh"]

