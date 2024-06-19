package main

import (
	"log"
)

func main() {
	distanceCalculator, err := NewDistanceCalculator()
	if err != nil {
		log.Fatal(err)
	}
	distanceCalculator.consumer.ConsumeData()
}

type DistanceCalculator struct {
	consumer DataConsumer
}

func NewDistanceCalculator() (*DistanceCalculator, error) {
	var (
		consumer DataConsumer
		addr     = "localhost:9092"
		topic    = "obudata"
		group    = "distance-calculator"
		err      error
	)
	consumer, err = NewKafkaConsumer(addr, topic, group)
	if err != nil {
		return nil, err
	}
	return &DistanceCalculator{
		consumer: consumer,
	}, nil
}
