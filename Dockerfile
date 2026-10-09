FROM golang:1.26 AS builder

RUN apt-get update && apt-get install -y --no-install-recommends \
        libvips-dev \
        pkg-config \
        && rm -rf /var/lib/apt/lists/

COPY . /src
WORKDIR /src

RUN GOPROXY=https://goproxy.cn make build

FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates \
        netbase \
        libvips42t64 \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/bin /app

WORKDIR /app

EXPOSE 8010
EXPOSE 9010
VOLUME /data/config

CMD ["./server", "-config", "/data/config"]