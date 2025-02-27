package main

import (
	"fmt"
	random "github.com/mazen160/go-random"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	"math/rand"
	"os"
	"strconv"
	"sync"
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
var message_each_batch = os.Getenv("MESSAGE_EACH_BATCH")
var rand_msg_size_in_mb = os.Getenv("RAND_MSG_SIZE_IN_MB")
var close_msg_log = os.Getenv("CLOSE_MSG_LOG")

var charset = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

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
	for {
		log.Println("Starting new publisher.........")
		publish()
		log.Println("Waiting to start new publisher.........")
		time.Sleep(5 * time.Second)
	}
}

// func submit(writer http.ResponseWriter, request *http.Request, p httprouter.Params) {
func publish() {
	// connect to db via amqp
	conn := connectToDB()
	if conn == nil {
		return
	}
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
			log.Printf("Failed to close db connection")
		}
	}(conn)

	// open a channel through db connection
	ch := getConnectionChannel(conn)
	if ch == nil {
		return
	}
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			log.Printf("Failed to open a channel through DB connection")
		}
	}(ch)

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
		log.Printf("%s: %s", "Failed to declare a queue", err)
		return
	}
	fmt.Printf("Declared a %s queue with name: %s\n", rabbit_queue_type, rabbit_queue)

	forever := make(chan bool)
	mainBlocker := make(chan bool)
	message_count := 0
	isPublishFailed := false

	tickerTime, err := strconv.Atoi(rabbit_publish_interval)
	if err != nil {
		return
	}
	ticker := time.NewTicker(time.Duration(tickerTime) * time.Second)
	defer ticker.Stop()
	go heartbeat(forever, ticker)

	go func() {
		for {
			select {
			case <-forever:
				var wg sync.WaitGroup
				message_total, _ := strconv.Atoi(message_each_batch)
				wg.Add(message_total)

				for i := 1; i <= message_total; i++ {
					go func() {
						defer wg.Done()
						message := ""
						msg_len, _ := strconv.Atoi(rand_msg_size_in_mb)
						if msg_len != 0 {
							message = generateRandomMessagesOfSizeInMB(msg_len)
						} else {
							t := time.Now()
							message = fmt.Sprintf("This is a message published at %d : %d : %d\n", t.Hour(), t.Minute(), t.Second())
						}

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
							log.Printf("%s: %s", "Failed to publish a message", err)
							isPublishFailed = true
							return
						}
						message_count++
						isPublishFailed = false
						val, _ := strconv.ParseBool(close_msg_log)
						if !val {
							fmt.Printf("Message Published: %s\n", message)
						}
						fmt.Printf("Total Message published: %v\n", message_count)
					}()
				}
			}
			if isPublishFailed {
				break
			}
		}
		mainBlocker <- true
	}()

	fmt.Println("Running Publisher...")
	<-mainBlocker
	fmt.Println("Stopping publisher...")
}

func heartbeat(ch chan bool, ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			ch <- true
		}
	}
}

func connectToDB() *amqp.Connection {
	if rabbit_vhost == "" {
		rabbit_vhost = "/"
	} else if rabbit_vhost != "/" {
		rabbit_vhost = "/" + rabbit_vhost
	}
	conn, err := amqp.Dial("amqp://" + rabbit_user + ":" + rabbit_password + "@" + rabbit_host + ":" + rabbit_port + rabbit_vhost)
	if err != nil {
		log.Printf("%s: %s", "Failed to connect to RabbitMQ", err)
		return nil
	}

	fmt.Println("Connected to RabbitMQ...")
	return conn
}

func getConnectionChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("%s: %s", "Failed to open a channel", err)
		return nil
	}
	fmt.Println("Opened channel to process AMQP messages")
	return ch
}

func generateRandomMessagesOfSizeInMB(sizeInMB int) string {
	randString, err := random.String(1024 * 1024 * sizeInMB)
	if err != nil {
		log.Println(err)
		return ""
	}
	return randString
}
