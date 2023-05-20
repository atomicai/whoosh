package handler

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func OptimalPath() {
	conn, err := amqp.Dial("amqp://rabbitmq:rabbitmq@localhost:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"DijkstraPathQuery", // name
		false,               // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	FailOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {

		log.Printf("Received a message: %s", msgs)

		q, err = ch.QueueDeclare(
			"DijkstraPathArgs", // name
			false,              // durable
			false,              // delete when unused
			false,              // exclusive
			false,              // no-wait
			nil,                // arguments
		)
		FailOnError(err, "Failed to declare a queue")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		body := "Hello World!"
		err = ch.PublishWithContext(ctx,
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		FailOnError(err, "Failed to publish a message")
		log.Printf(" [x] Sent %s\n", body)
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
