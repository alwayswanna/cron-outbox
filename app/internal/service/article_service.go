package service

import (
	"cron-outbox/internal/configuration"
	"cron-outbox/internal/db"
	"cron-outbox/internal/db/entity"
	"cron-outbox/internal/db/repository"
	"cron-outbox/internal/model"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

const ArticleTopicName = "article"

type ArticleService interface {
	ProceedSaveArticle(article model.ArticleRequest) error
	ProceedSendArticle()
}

type ArticleServiceImpl struct {
	properties              *configuration.Properties
	articleRepository       *repository.ArticleRepositoryImpl
	outboxMessageRepository *repository.OutboxMessageRepositoryImpl
	databaseConnection      *db.DatabaseConnection
	kafkaProducerService    *KafkaProducerService
}

func (a ArticleServiceImpl) ProceedSaveArticle(article model.ArticleRequest) error {
	var articleEntity = entity.Article{
		Title:       article.Title,
		Description: article.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	getDB := a.databaseConnection.GetDB()
	tx := getDB.Begin()

	savedEntity, err := a.articleRepository.Save(articleEntity)
	if err != nil {
		tx.Rollback()
		return err

	}

	outboxMessage, marshalErr := json.Marshal(savedEntity)
	if marshalErr != nil {
		tx.Rollback()
		return marshalErr
	}

	var outboxMessageToSave = entity.OutboxMessage{Message: string(outboxMessage)}
	_, err = a.outboxMessageRepository.Save(outboxMessageToSave)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (a ArticleServiceImpl) ProceedSendArticle() {
	articleWhereIsSentFalse, err := a.outboxMessageRepository.FindArticleByIsSent(false)
	if err != nil {
		log.Err(err).Msg("failed to find article where is sent false")
		return
	}

	if len(articleWhereIsSentFalse) == 0 {
		log.Info().Msg("there are no article to send")
		return
	} else {
		log.Info().Int("count", len(articleWhereIsSentFalse)).Msg("found articles to send")
	}

	/* check producer status, and attempt to recreate if closed */
	err = a.kafkaProducerService.PrepareToSendMessage(ArticleTopicName)
	if err != nil {
		log.Err(err).Msgf("producer for topic %s is closed or not exists", ArticleTopicName)
		return
	}

	/* process slice of articles */
	for _, article := range articleWhereIsSentFalse {
		err := a.databaseConnection.GetDB().Transaction(func(tx *gorm.DB) error {

			dbErr := tx.Model(&article).Updates(map[string]interface{}{
				"is_sent":    true,
				"updated_at": time.Now(),
			}).Error

			if dbErr != nil {
				return dbErr
			}

			err := a.kafkaProducerService.SendMessageToTopicInTransaction(
				article.Id.String(),
				article.Message,
				ArticleTopicName,
				tx,
			)

			if err != nil {
				tx.Rollback()
				return err
			}
			return nil
		})

		if err != nil {
			log.Err(err).Msg("failed to send message to topic")
		}
	}
}

func NewArticleService(
	properties *configuration.Properties,
	articleRepository *repository.ArticleRepositoryImpl,
	outboxMessageRepository *repository.OutboxMessageRepositoryImpl,
	databaseConnection *db.DatabaseConnection,
	kafkaProducerService *KafkaProducerService,
) *ArticleServiceImpl {
	return &ArticleServiceImpl{
		properties:              properties,
		articleRepository:       articleRepository,
		outboxMessageRepository: outboxMessageRepository,
		databaseConnection:      databaseConnection,
		kafkaProducerService:    kafkaProducerService,
	}
}
