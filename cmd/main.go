package main

import (
	"github.com/streadway/amqp"
	"log"
	"os"
	"sync"
)

type Feed struct {
	Name string
}

var wgSend sync.WaitGroup
var wgReceive sync.WaitGroup

func main() {
	StartFeed()
}

func StartFeed() {
	// Define RabbitMQ server URL.
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	connectRabbitMQ, err := amqp.Dial(amqpServerURL)

	if err != nil {
		panic(err)
	}

	defer connectRabbitMQ.Close()

	rabbitMQChannel, err2 := connectRabbitMQ.Channel()

	if err2 != nil {
		panic(err2)
	}

	defer rabbitMQChannel.Close()

	// Exchange
	err = rabbitMQChannel.ExchangeDeclare(
		"feeds",  // name
		"direct", // type
		false,    // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := rabbitMQChannel.QueueDeclare(
		"feed", // name
		false,  // durable
		false,  // delete when unused
		true,   // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = rabbitMQChannel.QueueBind(
		q.Name,     // queue name
		"NEW_FEED", // routing key
		"feeds",    // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	messages, err3 := rabbitMQChannel.Consume(
		q.Name, // queue name
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // arguments
	)
	if err3 != nil {
		log.Println(err)
	}

	// Build a welcome message.
	log.Println("Successfully connected to RabbitMQ")
	log.Println("Waiting for messages")

	// Make a channel to receive messages into infinite loop.
	forever := make(chan bool)

	go func() {
		for message := range messages {
			// For example, show received message in a console.
			log.Printf(" > Received message: %s\n", message.Body)
		}
	}()

	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
