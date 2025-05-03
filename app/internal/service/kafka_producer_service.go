package service

import (
	"context"
	"cron-outbox/internal/configuration"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
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

func (kafkaProducerService *KafkaProducerService) SendMessageToTopicInTransaction(
	messageID string,
	message string,
	topicName string,
	tx *gorm.DB,
) error {
	producer := kafkaProducerService.producers[topicName]
	if producer == nil || producer.IsClosed() {
		err := fmt.Errorf("producer for topic %s not found", topicName)
		log.Err(err)
		return err
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
		err := fmt.Errorf("failed to begin transaction: %s", err)
		log.Err(err)
		return err
	}

	deliveryChan := make(chan kafka.Event)
	err = producer.Produce(msg, deliveryChan)

	if err != nil {
		_ = producer.AbortTransaction(context.Background())
		close(deliveryChan)
		log.Err(err).Msg("failed to produce message to topic")
		return err
	}

	select {
	case e := <-deliveryChan:
		msg := e.(*kafka.Message)
		if msg.TopicPartition.Error != nil {
			_ = producer.AbortTransaction(context.Background())
			tx.Rollback()
			return msg.TopicPartition.Error
		}
		err = producer.CommitTransaction(context.Background())

		if err != nil {
			log.Err(err).Msg("failed to commit transaction")
			tx.Rollback()
			return err
		}

		return nil

	case <-time.After(10 * time.Second):
		_ = producer.AbortTransaction(context.Background())
		err := fmt.Errorf("kafka produce timeout")
		tx.Rollback()
		log.Err(err)
		return err
	}
}

func NewKafkaProducerService(this *configuration.Properties) *KafkaProducerService {
	producers := make(map[string]*kafka.Producer)
	txId, err := uuid.NewRandom()
	if err != nil {

	}
	for _, it := range this.KafkaProperties.Producers {
		producerConfigMap := &kafka.ConfigMap{
			BootstrapServersKey:  this.GetKafkaBootstrapServers(),
			EnableIdempotenceKey: it.EnableIdempotence,
			TransactionIDKey:     fmt.Sprintf("%s-%s", it.TransactionID, txId.String()),
			ClientIdKey:          it.ClientId,
			AcksKey:              it.Ack,
			RetriesKey:           it.Retries,
		}

		producer, err := kafka.NewProducer(producerConfigMap)
		if err != nil {
			log.Err(err).Msg("failed to create kafka producer")
		}

		err = producer.InitTransactions(context.Background())
		if err != nil {
			producer.Close()
			log.Err(err).Msg("failed to init transactions")
		}

		producers[it.TopicName] = producer
	}

	log.Info().Msg("kafka producer service successfully created")
	return &KafkaProducerService{properties: this, producers: producers}
}
