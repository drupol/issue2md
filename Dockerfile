# 使用 Go 语言官方镜像作为构建阶段
FROM golang:1.21 AS builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod文件
COPY go.mod ./

# 下载依赖
RUN go mod tidy

# 复制源代码
COPY . .

# 构建二进制文件
RUN CGO_ENABLED=0 GOOS=linux go build -o issue2mdweb .

# 使用轻量级的镜像作为生产阶段
FROM alpine:latest

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/issue2mdweb .

# 复制 web 目录下的 templates 和 static 子目录
COPY --from=builder /app/web/templates /root/web/templates
COPY --from=builder /app/web/static /root/web/static

# 暴露应用端口
EXPOSE 8080

# 指定容器启动时的命令
ENTRYPOINT ["./issue2mdweb"]
