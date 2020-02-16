package main

import (
	"log"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func sendMessage(ch *amqp.Channel, q amqp.Queue, message string) {
	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	log.Printf(" [x] Sent %s", message)
	failOnError(err, "Failed to publish a message")
}

func main() {
	conn, err := amqp.Dial("amqp://wbe7:mindbox123@34.89.173.76:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello2", // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue")

	for i := 0; i < 1000000; i++ {
		body := "Hello World!" + strconv.Itoa(i)
		sendMessage(ch, q, body)
		time.Sleep(500 * time.Microsecond)
	}
}
