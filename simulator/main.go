package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"simulator/entity"
	"simulator/queue"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

// init é executado por padrão antes da função main
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed loading .env:%s", err)
	}
}

func main() {
	// channels: nada mais é que um canal por onde será trafegado dados em go.
	// podemos usa-los para criar uma comunicação entre goroutinas
	// já que elas são executadas em threads diferentes
	inbox := make(chan []byte)
	ch := queue.Connect()           // canal de conexão com a fila
	queue.StartConsuming(inbox, ch) // iniciando o consumer

	// interamos nosso channel inbox, para buscar todos as de msgs que foram gravados nele
	for msg := range inbox {
		var order entity.Order

		// decodamos a msg em uma estrutura order
		if err := json.Unmarshal(msg, &order); err != nil {
			log.Printf("Failed unmarshal order:%s\n", err.Error())
		}
		log.Printf("Novo pedido realizado:%s\n", order.UUID)

		start(order, ch)
	}
}

func start(order entity.Order, ch *amqp.Channel) {
	go Worker(order, ch)
}

// Worker realiza a leitura da localização do destinatário
func Worker(order entity.Order, ch *amqp.Channel) {
	file, err := os.Open("destination/" + order.Destination + ".txt")

	if err != nil {
		log.Fatalf("Failed to open route %s: %s", order.Destination, err)
	}

	defer file.Close()

	s := bufio.NewScanner(file)
	for s.Scan() {
		coordinates := strings.Split(s.Text(), ",")
		destinationJSON := destination(order, coordinates[0], coordinates[1])

		time.Sleep(2 * time.Second)
		queue.Notify(string(destinationJSON), ch)
	}

	if s.Err() != nil {
		log.Fatalf("Failed buifio scanner: %s", s.Err())
	}

	destinationJSON := destination(order, "", "")
	queue.Notify(string(destinationJSON), ch)
}

//destination encoda a estrutura do destinatario em um json
func destination(order entity.Order, lat string, lgn string) []byte {
	dest := entity.Destination{
		Order:     order.UUID,
		Latitude:  lat,
		Longitude: lgn,
	}

	destinationJSON, _ := json.Marshal(dest)
	return destinationJSON
}
