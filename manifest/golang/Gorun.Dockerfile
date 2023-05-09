# syntax=docker/dockerfile:1.3

FROM 10.12.0.78:5000/cosmos/ci:golang-1.19 AS build

ARG flags

COPY . /src
WORKDIR /src

RUN --mount=type=cache,mode=0777,target=/root/.cache/go-build --mount=type=cache,mode=0777,target=/root/go/pkg go build -ldflags "-linkmode external -extldflags -static -s -w $flags" -o /bin/mesh


FROM --platform=amd64 10.12.0.78:5000/cosmos/ci:golang-1.19 AS compres

COPY --from=build /bin/mesh /bin/mesh
RUN upx -9 /bin/mesh


FROM scratch as export

COPY --from=compres /bin/mesh .
