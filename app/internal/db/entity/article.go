package entity

import (
	"github.com/google/uuid"
	"time"
)

type Article struct {
	Id          uuid.UUID `gorm:"id;type:uuid;default:gen_random_uuid();primary_key" json:"id"`
	Title       string    `gorm:"title;type:text;not null" json:"title"`
	Description string    `gorm:"description;type:text;not null" json:"description"`
	CreatedAt   time.Time `gorm:"created_at;type:timestamp;not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"updated_at;type:timestamp;not null" json:"updated_at"`
}

func (Article) TableName() string {
	return "article"
}
