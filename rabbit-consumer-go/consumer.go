package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
)

var rabbit_host = os.Getenv("RABBIT_HOST")
var rabbit_port = os.Getenv("RABBIT_PORT")
var rabbit_user = os.Getenv("RABBIT_USERNAME")
var rabbit_password = os.Getenv("RABBIT_PASSWORD")
var rabbit_queue = os.Getenv("RABBIT_QUEUE")
var rabbit_vhost = os.Getenv("RABBIT_VHOST")

func main() {
	consume()
}

func consume() {

	if rabbit_vhost == "" {
		rabbit_vhost = "/"
	} else if rabbit_vhost != "/" {
		rabbit_vhost = "/" + rabbit_vhost
	}
	conn, err := amqp.Dial("amqp://" + rabbit_user + ":" + rabbit_password + "@" + rabbit_host + ":" + rabbit_port + rabbit_vhost)

	if err != nil {
		log.Fatalf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	ch, err := conn.Channel()

	if err != nil {
		log.Fatalf("%s: %s", "Failed to open a channel", err)
	}

	fmt.Println("Channel established")

	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		rabbit_queue,    // queue
		"test-consumer", // consumer
		false,           // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)

	if err != nil {
		log.Fatalf("%s: %s", "Failed to register consumer", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			fmt.Printf("Received a message: %s\n", d.Body)

			d.Ack(false)
		}
	}()

	fmt.Println("Running...")
	<-forever
}
