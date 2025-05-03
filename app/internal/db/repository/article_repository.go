package repository

import (
	"cron-outbox/internal/db"
	"cron-outbox/internal/db/entity"
)

type ArticleRepository interface {
	Save(article entity.Article) (entity.Article, error)
}

type ArticleRepositoryImpl struct {
	db *db.DatabaseConnection
}

func (a *ArticleRepositoryImpl) Save(article entity.Article) (entity.Article, error) {
	gorm := a.db.GetDB()
	gorm.Create(&article)
	return article, nil
}

func NewArticleRepository(db *db.DatabaseConnection) *ArticleRepositoryImpl {
	return &ArticleRepositoryImpl{db: db}
}
