package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
	"strconv"
	"time"
)

var rabbit_host = os.Getenv("RABBIT_HOST")
var rabbit_port = os.Getenv("RABBIT_PORT")
var rabbit_user = os.Getenv("RABBIT_USERNAME")
var rabbit_password = os.Getenv("RABBIT_PASSWORD")
var rabbit_queue = os.Getenv("RABBIT_QUEUE")
var rabbit_queue_type = os.Getenv("RABBIT_QUEUE_TYPE")
var rabbit_publish_interval = os.Getenv("RABBIT_PUBLISH_INTERVAL")
var rabbit_vhost = os.Getenv("RABBIT_VHOST")

func main() {
	//
	//router := httprouter.New()
	//
	//router.POST("/publish/:message", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//	submit(w, r, p)
	//})
	//
	//fmt.Println("Running...")
	//log.Fatal(http.ListenAndServe(":80", router))
	publish()
}

// func submit(writer http.ResponseWriter, request *http.Request, p httprouter.Params) {
func publish() {

	if rabbit_vhost == "" {
		rabbit_vhost = "/"
	} else if rabbit_vhost != "/" {
		rabbit_vhost = "/" + rabbit_vhost
	}
	conn, err := amqp.Dial("amqp://" + rabbit_user + ":" + rabbit_password + "@" + rabbit_host + ":" + rabbit_port + rabbit_vhost)

	if err != nil {
		log.Fatalf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	fmt.Println("Connected to RabbitMQ...")

	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		log.Fatalf("%s: %s", "Failed to open a channel", err)
	}

	fmt.Println("Opened channel to process AMQP messages")

	defer ch.Close()

	args := amqp.Table{"x-queue-type": rabbit_queue_type}
	q, err := ch.QueueDeclare(
		rabbit_queue, // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		args,         // arguments
	)

	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare a queue", err)
	}

	fmt.Printf("Declared a %s queue with name: %s\n", rabbit_queue_type, rabbit_queue)

	forever := make(chan bool)
	mainBlocker := make(chan bool)

	go heartbeat(forever)

	go func() {
		for {
			select {
			case <-forever:
				t := time.Now()
				message := fmt.Sprintf("This is a message published at %d : %d : %d\n", t.Hour(), t.Minute(), t.Second())

				err = ch.Publish(
					"",     // exchange
					q.Name, // routing key
					false,  // mandatory
					false,  // immediate
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        []byte(message),
					})

				if err != nil {
					log.Fatalf("%s: %s", "Failed to publish a message", err)
				}

				fmt.Printf("Message Published: %s\n", message)
			}
		}
	}()

	fmt.Println("Running...")
	<-mainBlocker
}

func heartbeat(ch chan bool) {
	tickerTime, err := strconv.Atoi(rabbit_publish_interval)
	if err != nil {
		return
	}

	ticker := time.NewTicker(time.Duration(tickerTime) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			ch <- true
		}
	}
}
