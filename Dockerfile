FROM golang:1.17-alpine AS builder

WORKDIR /app

# 设置 GOPROXY 环境变量
ENV GOPROXY=https://goproxy.cn,direct

# 安装 git
RUN apk add --no-cache git

# 复制整个项目目录
COPY . .

# 初始化 go module，下载依赖，并生成 go.mod 和 go.sum
RUN go mod init shodan-proxy && go mod tidy

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o shodan-proxy ./cmd/server

FROM alpine:latest

# 添加作者标识
LABEL maintainer="LIUWEITAO <me@liuweitao.cn>"

WORKDIR /app

COPY --from=builder /app/shodan-proxy .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/public ./public
COPY --from=builder /app/config ./config

EXPOSE 8080

CMD ["./shodan-proxy"]
