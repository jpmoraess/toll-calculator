package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
	"github.com/jpmoraess/toll-calculator/common"
)

type DataConsumer interface {
	ConsumeData()
}

type KafkaConsumer struct {
	consumer sarama.ConsumerGroup
	addr     string
	topic    string
	group    string
}

func NewKafkaConsumer(addr, topic, group string) (DataConsumer, error) {
	brokers := []string{addr}
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()

	consumer, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		log.Fatalf("erro ao criar o consumer: %v", err)
		return nil, err
	}

	return &KafkaConsumer{
		consumer: consumer,
		addr:     addr,
		topic:    topic,
		group:    group,
	}, nil
}

func (c *KafkaConsumer) ConsumeData() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumerHandler := ConsumerGroupHandler{}

	go func() {
		for {
			if err := c.consumer.Consume(context.Background(), []string{c.topic}, consumerHandler); err != nil {
				log.Fatalf("erro ao consumir mensagens do kafka no tópico: %s", c.topic)
			}
		}
	}()

	<-signals
}

type ConsumerGroupHandler struct{}

// Setup é executado ao iniciar o Consumer Group
func (ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup é executado ao finalizar o Consumer Group
func (ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var message common.OBUData
		if err := json.Unmarshal(msg.Value, &message); err != nil {
			log.Printf("erro ao desserializar mensagem: %v", err)
			continue
		}
		log.Printf("mensagem recebida com sucesso: %+v", message)
	}
	return nil
}
