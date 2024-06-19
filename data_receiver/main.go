package main

import (
	"fmt"
	"log"
	"net/http"

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
	producer DataProducer
}

func NewDataReceiver() (*DataReceiver, error) {
	var (
		producer DataProducer
		topic    = "obudata"
		addr     = "localhost:9092"
		err      error
	)
	producer, err = NewKafkaProducer(addr, topic)
	if err != nil {
		return nil, err
	}
	producer = NewLogMiddleware(producer)
	return &DataReceiver{
		msgCh:    make(chan common.OBUData),
		producer: producer,
	}, nil
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

		if err := dr.producer.ProduceData(data); err != nil {
			fmt.Println("kafka producer error:", err)
		}
	}
}
