package service

import (
	"context"
	"cron-outbox/internal/configuration"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"slices"
)

const (
	BootstrapServersKey  = "bootstrap.servers"
	ClientIdKey          = "client.id"
	AcksKey              = "acks"
	RetriesKey           = "retries"
	EnableIdempotenceKey = "enable.idempotence"
	TransactionIDKey     = "transactional.id"
)

type KafkaProducerService struct {
	properties *configuration.Properties
	producers  map[string]*kafka.Producer
}

func (k *KafkaProducerService) PrepareToSendMessage(topicName string) error {
	if k.producers[topicName] != nil && k.producers[topicName].IsClosed() {
		producerConfigIdx := slices.IndexFunc(k.properties.KafkaProperties.Producers, func(producer configuration.Producers) bool {
			return producer.TopicName == topicName
		})
		producerProps := k.properties.KafkaProperties.Producers[producerConfigIdx]

		if &producerProps != nil {
			producer, err := _fromProducerProperties(&producerProps, k.properties.GetKafkaBootstrapServers())
			if err != nil {
				log.Err(err).Msgf("failed to recreate producer for [topic_name=%s]", topicName)
				return err
			}
			k.producers[topicName] = producer
		}
	}

	return nil
}

func (k *KafkaProducerService) SendMessageToTopicInTransaction(
	messageID string,
	message string,
	topicName string,
	tx *gorm.DB,
) error {
	producer := k.producers[topicName]
	if producer == nil || producer.IsClosed() {
		err := fmt.Errorf("producer for topic %s not found", topicName)
		log.Err(err)
		tx.Rollback()
		return err
	}

	onError := func(err error, errorMsg string) {
		log.Err(err).Msgf(errorMsg, topicName, message)
		tx.Rollback()
		producer.Close()
	}

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topicName,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(message),
		Headers: []kafka.Header{
			{Key: "message_id", Value: []byte(messageID)},
		},
		Opaque: tx,
	}

	err := producer.BeginTransaction()
	if err != nil {
		onError(err, "failed to begin transaction")
		return err
	}

	deliveryChan := make(chan kafka.Event)
	err = producer.Produce(msg, deliveryChan)

	if err != nil {
		onError(err, "failed to produce message")
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		onError(m.TopicPartition.Error, "failed to deliver message")
		return m.TopicPartition.Error
	} else {
		err := producer.CommitTransaction(context.Background())
		if err != nil {
			onError(err, "failed to commit transaction")
		}
		log.Info().Msg("message delivered successfully")
		close(deliveryChan)
		return nil
	}
}

func NewKafkaProducerService(this *configuration.Properties) *KafkaProducerService {
	producers := make(map[string]*kafka.Producer)

	for _, it := range this.KafkaProperties.Producers {
		producer, err := _fromProducerProperties(&it, this.GetKafkaBootstrapServers())

		if err != nil {
			log.Err(err).Msgf("Failed to create kafka producer")
			panic(err)
		}

		producers[it.TopicName] = producer
	}

	return &KafkaProducerService{properties: this, producers: producers}
}

func _fromProducerProperties(producerProperties *configuration.Producers, bootstrapAddress string) (*kafka.Producer, error) {
	txId, err := uuid.NewRandom()
	if err != nil {
		log.Err(err).Msg("failed to generate transaction id")
		return nil, err
	}

	producerConfigMap := &kafka.ConfigMap{
		BootstrapServersKey:  bootstrapAddress,
		EnableIdempotenceKey: producerProperties.EnableIdempotence,
		TransactionIDKey:     fmt.Sprintf("%s-%s", producerProperties.TransactionID, txId.String()),
		ClientIdKey:          producerProperties.ClientId,
		AcksKey:              producerProperties.Ack,
		RetriesKey:           producerProperties.Retries,
	}

	producer, err := kafka.NewProducer(producerConfigMap)
	if err != nil {
		log.Err(err).Msg("failed to create kafka producer")
		return nil, err
	}

	err = producer.InitTransactions(context.Background())
	if err != nil {
		producer.Close()
		log.Err(err).Msg("failed to init transactions")
		return nil, err
	}

	return producer, nil
}
