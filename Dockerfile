FROM rabbitmq:3-management AS rabbitmq
RUN rabbitmq-plugins enable rabbitmq_web_stomp


FROM golang:1.20.4-alpine AS goAPP

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./app ./

RUN go build -o /main

CMD [ "/main" ]