package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/IBM/sarama"
	"github.com/jpmoraess/toll-calculator/common"
	"github.com/jpmoraess/toll-calculator/distance_calculator/client"
	"github.com/sirupsen/logrus"
)

type DataConsumer interface {
	Consume() error
}

type KafkaConsumer struct {
	consumer         sarama.ConsumerGroup
	addr             string
	topic            string
	group            string
	service          CalculatorServicer
	aggregatorClient client.AggregatorClient
}

func NewKafkaConsumer(addr, topic, group string, service CalculatorServicer, aggregatorClient client.AggregatorClient) (DataConsumer, error) {
	brokers := []string{addr}
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()

	consumer, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		log.Fatalf("erro ao criar o consumer: %v", err)
		return nil, err
	}

	return &KafkaConsumer{
		consumer:         consumer,
		addr:             addr,
		topic:            topic,
		group:            group,
		service:          service,
		aggregatorClient: aggregatorClient,
	}, nil
}

func (c *KafkaConsumer) Consume() error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumerHandler := &ConsumerGroupHandler{
		consumer: c,
	}

	go func() {
		for {
			if err := c.consumer.Consume(context.Background(), []string{c.topic}, consumerHandler); err != nil {
				log.Fatalf("erro ao consumir mensagens do kafka no tópico: %s", c.topic)
			}
		}
	}()

	<-signals

	return nil
}

type ConsumerGroupHandler struct {
	consumer *KafkaConsumer
}

// Setup é executado ao iniciar o Consumer Group
func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup é executado ao finalizar o Consumer Group
func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var message common.OBUData
		if err := json.Unmarshal(msg.Value, &message); err != nil {
			log.Printf("erro ao desserializar mensagem: %v", err)
			continue
		}
		distance, err := h.consumer.service.CalculateDistance(message)
		if err != nil {
			return err
		}

		req := common.Distance{
			OBUID: message.OBUID,
			Value: distance,
			Unix:  time.Now().UnixNano(),
		}
		if err := h.consumer.aggregatorClient.AggregateInvoice(context.Background(), req); err != nil {
			logrus.Errorf("aggregate error: %v", err)
			continue
		}

	}
	return nil
}
