FROM golang:latest
# RUN go build -v
WORKDIR /home/go/release/app
COPY . .
LABEL Name=tcpserver Version=0.0.1
EXPOSE 12336
CMD ["./server"]
# 构建镜像  docker build .
# 标签镜像  docker tag [imageid]
# 运行容器  docker run -d -p 12336:12336 --name tcpserver tcpserver
# 查看容器日志 docker logs tcpserver