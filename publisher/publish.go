package main

import (
	"io/ioutil"
	"log"
	"os/exec"
	"time"

	"github.com/streadway/amqp"
)

func main() {

	// New RabbitMQ AMQP connection
	conn, err := amqp.Dial("amqp://user:bitnami@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	log.Println("Connected to RabbitMQ server")
	defer conn.Close()

	// Create a new Channel
	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	log.Println("Successfully created a new channel")
	defer channel.Close()

	// Declare a new Queue
	queue, err := channel.QueueDeclare(
		"log-queue", // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	log.Println("Created a new Queue: `log-queue`")
	failOnError(err, "Failed to declare a queue")

	// Starting Log stream
	log.Printf(" [*] Sending messages. To exit press CTRL+C")

	forever := make(chan bool)

	// Send messages infinitely
	go func() {
		for {
			cmd := exec.Command("flog", "-n", "1", "-f", "apache_combined")
			out, err := cmd.StdoutPipe()
			if err != nil {
				panic(err)
			}
			cmd.Start()
			messageBody, _ := ioutil.ReadAll(out)

			// Publish log to queue
			err = channel.Publish(
				"",         // exchange
				queue.Name, // routing key
				false,      // mandatory
				false,      // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(messageBody),
				})
			failOnError(err, "Failed to publish a message")

			// Sleep for 500ms before next log
			time.Sleep(time.Millisecond * 400)
		}

	}()
	<-forever
}

func failOnError(err error, msg string) bool {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic("Shutting down")
	}
	return true
}
