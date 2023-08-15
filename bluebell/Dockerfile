FROM golang:alpine AS builder

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 移动到工作目录：/build
WORKDIR /build

# 将代码复制到容器中
COPY . .

# 将我们的代码编译成二进制可执行文件 app
RUN go build -o bluebull_app .

###################
# 接下来创建一个小镜像
###################
FROM scratch

# 从builder镜像中把静态文件拷贝到当前目录
COPY ./templates /templates
COPY ./static /static

# 从builder镜像中把配置文件拷贝到当前目录
COPY ./conf /conf

# 从builder镜像中把/dist/app 拷贝到当前目录
COPY --from=builder /build/bluebull_app /

# 声明服务端口 不起作用
EXPOSE 8081

# 需要运行的命令
ENTRYPOINT ["/bluebull_app", "conf/config.yaml"]