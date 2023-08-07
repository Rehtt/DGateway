FROM golang:1.20.7-alpine3.18 as builder

COPY . /build
WORKDIR /build

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
    apk add --no-cache git
RUN export GOPROXY=https://goproxy.cn && \
    VERSION=$(git describe --tags --always) && \
	go mod tidy && \
    go build -o run -tags=jsoniter -ldflags "-X main.Version=$VERSION" .

FROM alpine:3.18

COPY --from=builder /build/run /app/run

WORKDIR /app

CMD ["/app/run"]