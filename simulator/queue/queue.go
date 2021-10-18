package queue

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

// Connect cria a conexão com o RabbitMQ
func Connect() *amqp.Channel {
	conn, err := amqp.Dial(getConnectString())
	failOnError(err, "Failed to connect RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open channel")

	return ch
}

func StartConsuming(inbox chan []byte, ch *amqp.Channel) {
	// criação da fila
	q, err := ch.QueueDeclare(os.Getenv("RABBITMQ_CONSUMER_QUEUE"), true, false, false, false, nil)
	failOnError(err, "Failed to create queue")

	// registrando o consume que vai ler a fila
	msgs, err := ch.Consume(q.Name, "go-worker-simulator", true, false, false, false, nil)
	failOnError(err, "Failed register consumer")

	go func() {
		//realiza a leitura das mensagens e envia para dentro de um channel do go
		for msg := range msgs {
			inbox <- []byte(msg.Body)
		}
		close(inbox)
	}()
}

// Notify realiza a publicação na msg em um exchange
func Notify(payload string, ch *amqp.Channel) {
	err := ch.Publish(
		os.Getenv("RABBITMQ_DESTINATION"),
		os.Getenv("RABBITMQ_DESTINATION_ROUTING_KEY"),
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(payload),
		},
	)

	failOnError(err, "Failed to publish message")
	log.Printf("Message send: %s\n", payload)
}

// getConnectString retorna a string de connexão com o RabbitMQ
func getConnectString() string {
	USER := os.Getenv("RABBITMQ_DEFAULT_USER")
	PASS := os.Getenv("RABBITMQ_DEFAULT_PASS")
	HOST := os.Getenv("RABBITMQ_DEFAULT_HOST")
	PORT := os.Getenv("RABBITMQ_DEFAULT_PORT")
	VHOST := os.Getenv("RABBITMQ_DEFAULT_VHOST")

	// dsn := "amqp://"+USER+":"+PASS+"@"+HOST+":"+PORT+VHOST
	return fmt.Sprintf("amqp://%s:%s@%s:%s%s", USER, PASS, HOST, PORT, VHOST)
}
