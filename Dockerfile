FROM golang:1.20.7-alpine3.18 as builder

COPY . /build
WORKDIR /build

RUN git describe --tags --always && \
	go mod tidy && \
    go build -o run -tags=jsoniter -ldflags "-X main.Version=$(VERSION)" .

FROM alpine:3.18

COPY --from=builder /build/run /app

WORKDIR /app

CMD ["/app/run"]