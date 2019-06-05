FROM golang:1.12.5-alpine

COPY . /app

WORKDIR /app

ENV GO111MODULE=on

RUN go build -mod=vendor -o /usr/bin/app

EXPOSE 8080

ENTRYPOINT ["app"]
