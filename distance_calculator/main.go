package main

import (
	"log"
)

// Transport (HTTP, GRPC, KAFKA) -> attach business logic to this transport

func main() {
	var (
		err     error
		service CalculatorServicer
	)
	service = NewCalculatorService()
	service = NewLogMiddleware(service)
	distanceCalculator, err := NewDistanceCalculator(service)
	if err != nil {
		log.Fatal(err)
	}
	distanceCalculator.consumer.Consume()
}

type DistanceCalculator struct {
	consumer DataConsumer
}

func NewDistanceCalculator(service CalculatorServicer) (*DistanceCalculator, error) {
	var (
		consumer DataConsumer
		addr     = "localhost:9092"
		topic    = "obudata"
		group    = "distance-calculator"
		err      error
	)
	consumer, err = NewKafkaConsumer(addr, topic, group, service)
	if err != nil {
		return nil, err
	}
	return &DistanceCalculator{
		consumer: consumer,
	}, nil
}
