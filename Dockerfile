FROM rabbitmq:3-management AS rabbitmq
RUN rabbitmq-plugins enable rabbitmq_web_stomp


FROM golang:1.20.4-alpine AS goAPP

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN go build -o /app .

ENTRYPOINT [ "/app" ]