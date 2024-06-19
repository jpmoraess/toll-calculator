package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/gorilla/websocket"
	"github.com/jpmoraess/toll-calculator/common"
)

func main() {
	receiver, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/ws", receiver.handleWS)
	http.ListenAndServe(":3000", nil)
}

type DataReceiver struct {
	msgCh    chan common.OBUData
	conn     *websocket.Conn
	producer sarama.SyncProducer
}

func NewDataReceiver() (*DataReceiver, error) {
	brokers := []string{"localhost:9092"}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("erro ao criar o producer: %v", err)
	}

	return &DataReceiver{
		msgCh:    make(chan common.OBUData),
		producer: producer,
	}, nil
}

func (dr *DataReceiver) produceData(data *common.OBUData) error {
	topic := "obudata"

	b, err := json.Marshal(&data)
	if err != nil {
		log.Fatalf("erro ao serializar a mensagem: %v", &data)
	}

	message := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder("Key-A"),
		Value: sarama.ByteEncoder(b),
	}

	partition, offset, err := dr.producer.SendMessage(message)
	if err != nil {
		log.Fatalf("erro ao enviar mensagem: %v", err)
	}
	log.Printf("mensagem enviada com sucesso, partition: %d, offset: %d", partition, offset)
	return nil
}

func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn
	go dr.receiverLoop()
}

func (dr *DataReceiver) receiverLoop() {
	fmt.Println("new OBU connected successfully")
	for {
		var data common.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read websocket error: ", err)
			continue
		}
		if err := dr.produceData(&data); err != nil {
			fmt.Println("kafka produce error:", err)
		}
	}
}
