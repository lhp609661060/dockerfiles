FROM golang:latest

MAINTAINER lhp "609661060@qq.com"

RUN go get github.com/gin-gonic/gin
RUN go get github.com/gomodule/redigo/redis
WORKDIR /go/src/app
COPY ./app /go/src/app
RUN go build -o main .

ENTRYPOINT ["/go/src/app/main"]

EXPOSE 80