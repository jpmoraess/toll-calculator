package main

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/jpmoraess/toll-calculator/common"
)

type DataProducer interface {
	ProduceData(common.OBUData) error
}

type KafkaProducer struct {
	producer sarama.SyncProducer
	addr     string
	topic    string
}

func NewKafkaProducer(addr, topic string) (DataProducer, error) {
	brokers := []string{addr}
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("erro ao criar o producer: %v", err)
		return nil, err
	}
	return &KafkaProducer{
		producer: producer,
		addr:     addr,
		topic:    topic,
	}, nil
}

func (p *KafkaProducer) ProduceData(data common.OBUData) error {
	b, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("erro ao serializar a mensagem: %v", err)
		return err
	}

	message := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(b),
	}

	partition, offset, err := p.producer.SendMessage(message)
	if err != nil {
		log.Fatalf("erro ao enviar mensagem: %v", err)
		return err
	}
	log.Printf("mensagem enviada com sucesso, partition: %d e offset: %d", partition, offset)
	return nil
}
