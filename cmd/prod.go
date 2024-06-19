package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

func main3() {
	// Configurações do consumidor Kafka
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Crie um novo consumidor Kafka
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Falha ao iniciar o consumidor Sarama: %v", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalf("Erro ao fechar o consumidor: %v", err)
		}
	}()

	// Capture sinais do sistema operacional para lidar com a finalização graciosa
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("Sinal de desligamento recebido, saindo...")
		os.Exit(0)
	}()

	// Crie um consumidor de partição para o tópico "your-topic"
	partitionConsumer, err := consumer.ConsumePartition("your-topic", 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("Falha ao criar consumidor de partição: %v", err)
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalf("Erro ao fechar consumidor de partição: %v", err)
		}
	}()

	// Processar mensagens
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Mensagem recebida: Partition %d, Offset %d: %s\n", msg.Partition, msg.Offset, string(msg.Value))
		case err := <-partitionConsumer.Errors():
			log.Printf("Erro ao consumir mensagem: %v\n", err)
		}
	}
}
