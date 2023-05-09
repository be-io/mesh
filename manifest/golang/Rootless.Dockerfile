# syntax=docker/dockerfile:1.3

FROM 10.12.0.78:5000/middleware-release-1.5.7/gaia-mesh:1.5.7-beta-d33c809

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache libcap
RUN setcap CAP_NET_BIND_SERVICE=+eip /bin/mesh
RUN addgroup -S app --gid 1001 && adduser -S app -G app -u 1001
RUN chown -R app /usr/local
RUN mkdir -p /var/log/be && chown -R app /var/log/be
RUN mkdir -p /mesh && chown -R app /mesh
WORKDIR /home/app
USER app