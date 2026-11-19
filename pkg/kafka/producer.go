package kafka

import (
	"encoding/json"
	"fmt"
	"time"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	ConfigMap *ckafka.ConfigMap
}

func NewKafkaProducer(configMap *ckafka.ConfigMap) *Producer {
	return &Producer{ConfigMap: configMap}
}

func (p *Producer) Publish(msg interface{}, key []byte, topic string) error {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	producer, err := ckafka.NewProducer(p.ConfigMap)
	if err != nil {
		return err
	}
	defer producer.Close()

	message := &ckafka.Message{
		TopicPartition: ckafka.TopicPartition{Topic: &topic, Partition: ckafka.PartitionAny},
		Value:          msgJSON,
		Key:            key,
	}
	deliveryChan := make(chan ckafka.Event, 1)
	err = producer.Produce(message, deliveryChan)
	if err != nil {
		return err
	}

	select {
	case event := <-deliveryChan:
		deliveryReport, ok := event.(*ckafka.Message)
		if !ok {
			return fmt.Errorf("unexpected Kafka delivery event %T", event)
		}
		return deliveryReport.TopicPartition.Error
	case <-time.After(10 * time.Second):
		return fmt.Errorf("timed out waiting for Kafka delivery to topic %q", topic)
	}
}
