FROM golang:1.17.6-alpine3.15 as builder

MAINTAINER sxueck sxuecks@gmail.com

COPY / /app
WORKDIR /app

RUN echo "https://mirrors.aliyun.com/alpine/v3.15/main" > /etc/apk/repositories
RUN echo "https://mirrors.aliyun.com/alpine/v3.15/community" >> /etc/apk/repositories

RUN apk add gcc build-base \
    && GO111MODULE=on GOPROXY="https://goproxy.cn" GOOS=linux go build -a -o /go/bin/starry-night .

FROM alpine:3.15

EXPOSE 80

RUN echo "https://mirrors.aliyun.com/alpine/v3.15/main" > /etc/apk/repositories
RUN echo "https://mirrors.aliyun.com/alpine/v3.15/community" >> /etc/apk/repositories

RUN apk add -U tzdata \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && apk add ca-certificates \
    && rm -rf /var/cache/apk/*
COPY --from=builder /go/bin/starry-night .

COPY public public/
COPY storage.db databases/

ENV DB_NAME "databases/storage.db"

ENV BASE_PATH /

CMD ["/starry-night","start"]
