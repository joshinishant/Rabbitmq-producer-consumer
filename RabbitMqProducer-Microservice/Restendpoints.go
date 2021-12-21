package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", defaultAction)
	router.HandleFunc("/Message", postMessage).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", router))
}

func postMessage(resp http.ResponseWriter, req *http.Request) {

	data, err := ioutil.ReadAll(req.Body)
	logError(err)

	conn, conn_err := amqp.Dial("amqp://guest:guest@local-rabbitmq:5672")
	defer conn.Close()
	logError(conn_err)

	ch, ch_err := conn.Channel()
	defer ch.Close()
	logError(ch_err)

	q, q_err := ch.QueueDeclare("INBOUND", true, false, false, false, nil)
	logError(q_err)

	ch.Publish("", q.Name, false, false, amqp.Publishing{Body: data, ContentType: "text/plain"})
	resp.WriteHeader(http.StatusAccepted)
	resp.Write([]byte("Message written to queue"))
}

func defaultAction(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusFound)
	resp.Write([]byte("Current enpoint is incorrect. Please send use correct url"))
}

func logError(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
