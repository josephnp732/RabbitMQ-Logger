package main

import (
	"log"

	"github.com/streadway/amqp"
)

const routingKey string = "logs-routing-key"
const exchangeName string = "logs-exchange"

func main() {
	conn, err := amqp.Dial("amqp://user:bitnami@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Create Channel
	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	// Declaring Queue here as well
	// Because we might start the consumer before the publisher,
	// we want to make sure the queue exists before we try to consume messages from it.
	queue, err := channel.QueueDeclare(
		"log-queue", // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Create new Exchange (Direct Exchange)
	err = channel.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	// Bind the Exchange and Queue
	err = channel.QueueBind(
		queue.Name,   // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil,
	)

	// Consume message
	msgs, err := channel.Consume(
		queue.Name,     // queue
		"log-consumer", // consumer name
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(string(d.Body))
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
