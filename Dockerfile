# syntax=docker/dockerfile:1.3

FROM 10.12.0.78:5000/cosmos/ci:golang-1.19 AS build

ARG flags

COPY . /src
WORKDIR /src

RUN --mount=type=cache,mode=0777,target=/root/.cache/go-build --mount=type=cache,mode=0777,target=/root/go/pkg go build -ldflags "-linkmode external -extldflags -static -s -w $flags" -o /bin/mesh


FROM --platform=amd64 10.12.0.78:5000/cosmos/ci:golang-1.19 AS compres

COPY --from=build /bin/mesh /bin/mesh
RUN upx -9 /bin/mesh


FROM alpine

MAINTAINER coyzeng@gmail.com

COPY --from=compres /bin/mesh /bin/mesh
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /src/promtail.yaml /etc/promtail/promtail.yaml

ENV TZ=Asia/Shanghai
ENV LANG=en_US.UTF-8

ENTRYPOINT ["/bin/mesh", "start", "--server"]
