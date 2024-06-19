package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jpmoraess/toll-calculator/common"
)

func main() {
	receiver := NewDataReceiver()
	http.HandleFunc("/ws", receiver.handleWS)
	http.ListenAndServe(":3000", nil)
}

type DataReceiver struct {
	msgCh chan common.OBUData
	conn  *websocket.Conn
}

func NewDataReceiver() *DataReceiver {
	return &DataReceiver{
		msgCh: make(chan common.OBUData),
	}
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
		fmt.Printf("received OBU data from ID: %v :: lat: %.2f :: long: %.2f\n", data.OBUID, data.Lat, data.Long)
		dr.msgCh <- data
	}
}
