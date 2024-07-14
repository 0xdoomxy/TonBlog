# 使用官方Golang镜像作为构建环境
FROM golang:1.22.1 as builder

# 设置工作目录
WORKDIR /app

# 复制go模块和依赖文件
COPY go.mod go.sum ./

RUN go env -w GO111MODULE=on


RUN go env -w GOPROXY=https://goproxy.cn,direct
# 下载依赖
RUN go mod download

# 复制源代码
COPY . .


# 构建应用程序
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o blog ./main.go

# 使用alpine作为最小运行时容器
FROM alpine



# 从构建器阶段复制编译好的应用程序
COPY --from=builder /app/blog .

# 复制配置文件
COPY --from=builder /app/config/*.yaml ./config/

EXPOSE 8080
# 运行应用程序
CMD ["./blog"]