package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/atomicai/whoosh/internal/models"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func OptimalPath() {
	conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	qReq, err := ch.QueueDeclare(
		"DijkstraPathQuery", // name
		false,               // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	qRes, err := ch.QueueDeclare(
		"DijkstraPathResponse", // name
		false,                  // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		qReq.Name, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		for d := range msgs {
			bodyToJson := models.PathQuery{}

			err := json.Unmarshal(d.Body, &bodyToJson)
			if err != nil {
				log.Fatal(err)
			}

			result1 := Dijkstra(&bodyToJson)
			//result2 := AStar(&bodyToJson)
			fmt.Printf("result1: %+v\n", result1)
			//fmt.Printf("result2: %+v\n", result2)

			err = ch.PublishWithContext(ctx,
				"",        // exchange
				qRes.Name, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(d.Body),
				})

		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
