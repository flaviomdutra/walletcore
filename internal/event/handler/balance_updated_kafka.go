package handler

import (
	"log"
	"sync"

	"github.com/flaviomdutra/walletcore/pkg/events"
	"github.com/flaviomdutra/walletcore/pkg/kafka"
)

type BalanceUpdatedKafkaHandler struct {
	Kafka *kafka.Producer
}

func NewBalanceUpdatedKafkaHandler(kafka *kafka.Producer) *BalanceUpdatedKafkaHandler {
	return &BalanceUpdatedKafkaHandler{
		Kafka: kafka,
	}
}

func (h *BalanceUpdatedKafkaHandler) Handle(message events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	if err := h.Kafka.Publish(message, nil, "balances"); err != nil {
		log.Printf("publish BalanceUpdated event to Kafka: %v", err)
	}
}
