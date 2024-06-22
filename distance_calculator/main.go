package main

import (
	"log"

	"github.com/jpmoraess/toll-calculator/distance_calculator/client"
)

// Transport (HTTP, GRPC, KAFKA) -> attach business logic to this transport

func main() {
	var (
		err              error
		service          CalculatorServicer
		aggregatorClient *client.AggregatorClient
	)
	aggregatorClient = client.NewAggregatorClient("http://localhost:3001/aggregate")
	service = NewCalculatorService()
	service = NewLogMiddleware(service)
	distanceCalculator, err := NewDistanceCalculator(service, aggregatorClient)
	if err != nil {
		log.Fatal(err)
	}
	distanceCalculator.consumer.Consume()
}

type DistanceCalculator struct {
	consumer DataConsumer
}

func NewDistanceCalculator(service CalculatorServicer, aggregatorClient *client.AggregatorClient) (*DistanceCalculator, error) {
	var (
		consumer DataConsumer
		addr     = "localhost:9092"
		topic    = "obudata"
		group    = "distance-calculator"
		err      error
	)
	consumer, err = NewKafkaConsumer(addr, topic, group, service, aggregatorClient)
	if err != nil {
		return nil, err
	}
	return &DistanceCalculator{
		consumer: consumer,
	}, nil
}
