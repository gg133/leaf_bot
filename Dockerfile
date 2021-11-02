FROM golang:1.17-alpine3.14 AS builder

COPY . /github.com/yalagtyarzh/leaf_bot/
WORKDIR /github.com/yalagtyarzh/leaf_bot/

RUN go mod download
RUN go build -o ./bin/bot main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/yalagtyarzh/leaf_bot/bin/bot .
COPY --from=0 /github.com/yalagtyarzh/leaf_bot/configs configs/

EXPOSE 80

CMD ["./bot"]