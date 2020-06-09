package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/streadway/amqp"
)

// number of publishes
var pubCount uint64

const routingKey string = "logs-routing-key"
const exchangeName string = "logs-exchange"

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

	// Create new Exchange (Direct Exchange)
	err = channel.ExchangeDeclare(
		exchangeName,        // name
		amqp.ExchangeDirect, // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	failOnError(err, "Failed to declare an exchange")

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

	// Bind the Exchange and Queue
	err = channel.QueueBind(
		queue.Name,   // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil,
	)

	// Starting Log stream
	log.Printf(" [*] Sending messages. To exit press CTRL+C")

	forever := make(chan bool)

	// Setup our Ctrl+C (Exit) handler
	SetupCloseHandler()

	// Send messages infinitely
	go func() {
		for {
			cmd := exec.Command("flog", "-n", "1", "-f", "apache_combined")
			// cmd := exec.Command("flog", "-n", "1", "-f", "json")
			out, err := cmd.StdoutPipe()
			if err != nil {
				panic(err)
			}
			cmd.Start()
			messageBody, _ := ioutil.ReadAll(out)

			// Publish log to queue
			err = channel.Publish(
				"logs-exchange", // exchange
				routingKey,      // routing key
				false,           // mandatory
				false,           // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(messageBody),
				})
			if !failOnError(err, "Failed to publish a message") {
				atomic.AddUint64(&pubCount, 1)
			}

			// Sleep for 500ms before next log
			time.Sleep(time.Millisecond * 500)
		}
	}()
	<-forever
}

func failOnError(err error, msg string) bool {
	if err != nil {
		defer log.Fatalf("%s: %s", msg, err)
	}
	return false
}

// SetupCloseHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an exit interrupt from the OS. We then handle this by calling
// our clean up procedure and exiting the program.
func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Exiting Program")
		log.Printf("Number of logs published: %d", pubCount)
		os.Exit(0)
	}()
}
