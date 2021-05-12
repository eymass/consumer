package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
	"strconv"
	"sync"
)

type Feed struct {
	Name string
}

var wg sync.WaitGroup

func main() {
	//StartFeed()
	TestFeed()
}

func TestFeed() {

	feedsChannel := make(chan Feed)
	fakeChan := make(chan Feed)

	go func() {
		for i := 0; i < 5; i++ {
			newFeed := Feed{
				Name: "Hello " + strconv.Itoa(i),
			}
			go sendMessage(feedsChannel, newFeed)
		}
	}()

	for i := 0; i < 5; i++ {
		go handleMessage(feedsChannel)
	}

	<-fakeChan
}

func sendMessage(feedsChannel chan Feed, msg Feed) {
	wg.Add(1)
	fmt.Println("sending, %s", msg.Name)
	feedsChannel <- msg
	defer wg.Done()
}

func handleMessage(feedsChannel chan Feed) {
	wg.Add(1)
	message2 := <-feedsChannel
	fmt.Println(message2.Name)
	defer wg.Done()
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

	messages, err3 := rabbitMQChannel.Consume(
		"Feed", // queue name
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
