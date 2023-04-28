# syntax=docker/dockerfile:1.3

FROM golang:1.20.7-alpine AS build

ARG flags

COPY . /src
WORKDIR /src

RUN go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,https://goproxy.cn,https://goproxy.io,https://proxy.golang.org,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add gcc g++ make git wget tzdata upx libpcap-dev openssh
RUN --mount=type=cache,mode=0777,target=/root/.cache/go-build --mount=type=cache,mode=0777,target=/root/go/pkg go build -ldflags "-linkmode external -extldflags -static -s -w $flags" -o /bin/mesh
RUN upx -9 /bin/mesh

FROM alpine

MAINTAINER coyzeng@gmail.com

COPY --from=build /bin/mesh /bin/mesh
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo

ENV TZ=Asia/Shanghai
ENV LANG=en_US.UTF-8

ENTRYPOINT ["/bin/mesh", "start", "--server"]
