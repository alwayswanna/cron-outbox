package repository

import (
	"cron-outbox/internal/db"
	"cron-outbox/internal/db/entity"
)

type OutboxMessageRepository interface {
	Save(message entity.OutboxMessage) (entity.OutboxMessage, error)
	FindArticleByIsSent(isSent bool) ([]entity.OutboxMessage, error)
	Update(message entity.OutboxMessage) (entity.OutboxMessage, error)
}

type OutboxMessageRepositoryImpl struct {
	db *db.DatabaseConnection
}

func (r *OutboxMessageRepositoryImpl) Save(message entity.OutboxMessage) (entity.OutboxMessage, error) {
	gorm := r.db.GetDB()
	gorm.Create(&message)
	return message, nil
}

func (r *OutboxMessageRepositoryImpl) FindArticleByIsSent(isSent bool) ([]entity.OutboxMessage, error) {
	var articles []entity.OutboxMessage

	gorm := r.db.GetDB()
	gorm.Find(&articles, "is_sent = ?", isSent)

	return articles, nil
}

func (r *OutboxMessageRepositoryImpl) Update(message entity.OutboxMessage) (entity.OutboxMessage, error) {
	gorm := r.db.GetDB()
	gorm.Save(&message)
	return message, nil
}

func NewOutboxMessageRepository(db *db.DatabaseConnection) *OutboxMessageRepositoryImpl {
	return &OutboxMessageRepositoryImpl{db: db}
}
