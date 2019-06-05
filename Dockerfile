FROM golang:alpine

COPY . /app

WORKDIR /app

ENV GO111MODULE=on

RUN go install -mod=vendor

EXPOSE 8080

ENTRYPOINT ["app"]
