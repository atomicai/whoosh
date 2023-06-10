FROM rabbitmq:3-management AS rabbitmq
RUN apt-get update && apt-get install -y wget curl
RUN rabbitmq-plugins enable rabbitmq_web_stomp

FROM golang:latest AS golang

WORKDIR /app

ADD go.mod .

COPY . .

RUN go build -o whoosh.exe ./cmd/app/main.go

CMD ["./whoosh.exe"]