package main

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jpmoraess/toll-calculator/common"
)

const websocketEndpoint = "ws://localhost:3000/ws"

func generateCoord() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()
	return n + f
}

func genLocation() (float64, float64) {
	return generateCoord(), generateCoord()
}

func generateIDS(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(math.MaxInt)
	}
	return ids
}

func main() {
	ids := generateIDS(20)
	conn, _, err := websocket.DefaultDialer.Dial(websocketEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		for i := 0; i < len(ids); i++ {
			lat, long := genLocation()
			data := common.OBUData{
				OBUID: ids[i],
				Lat:   lat,
				Long:  long,
			}
			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(time.Second)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
