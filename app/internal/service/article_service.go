package service

import (
	"cron-outbox/internal/configuration"
	"cron-outbox/internal/db/repository"
	"cron-outbox/internal/model"
)

type ArticleService interface {
	ProceedSaveArticle(article model.ArticleRequest) error
	ProceedSendArticle(article model.ArticleRequest) error
}

type ArticleServiceImpl struct {
	properties        *configuration.Properties
	articleRepository *repository.ArticleRepositoryImpl
}

func (a ArticleServiceImpl) ProceedSaveArticle(article model.ArticleRequest) error {
	//TODO implement me
	panic("implement me")
}

func (a ArticleServiceImpl) ProceedSendArticle(article model.ArticleRequest) error {
	//TODO implement me
	panic("implement me")
}

func NewArticleService(properties *configuration.Properties, articleRepository *repository.ArticleRepositoryImpl) ArticleService {
	return &ArticleServiceImpl{properties: properties, articleRepository: articleRepository}
}
