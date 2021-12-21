package main

import (
	"log"

	"github.com/streadway/amqp"
)

func main() {
	conn, conn_err := amqp.Dial("amqp://guest:guest@local-rabbitmq:5672")
	defer conn.Close()
	logError(conn_err)
	log.Println("Connection Established")

	ch, ch_err := conn.Channel()
	defer ch.Close()
	logError(ch_err)
	log.Println("Channel Created")

	q, q_err := ch.QueueDeclare("INBOUND", true, false, false, false, nil)
	logError(q_err)

	msgs, con_err := ch.Consume(q.Name, "", false, false, false, false, nil)
	logError(con_err)

	blockingChnl := make(chan bool)

	for msg := range msgs {
		log.Println("Message received - ", string(msg.Body))
		ack_err := ch.Ack(msg.DeliveryTag, false)
		logError(ack_err)
		log.Println("Message acknowledged now")
	}

	log.Println("Waiting for messages...")

	<-blockingChnl

}

func logError(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
