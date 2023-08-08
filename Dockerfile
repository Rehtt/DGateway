FROM golang:1.20.7-alpine3.18 as builder

COPY . /build
WORKDIR /build

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
    apk add --no-cache git make
RUN export GOPROXY=https://goproxy.cn && \
    VERSION=$(git describe --tags --always) && \
	make build

FROM alpine:3.18

COPY --from=builder /build/bin/run /app/run

WORKDIR /app

CMD ["/app/run"]