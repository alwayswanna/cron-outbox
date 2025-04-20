package entity

import (
	"github.com/google/uuid"
	"time"
)

type Article struct {
	Id          uuid.UUID `gorm:"id;type:uuid;default:uuid_generate_v4();primary_key"`
	Title       string    `gorm:"title;type:text;not null"`
	Description string    `gorm:"description;type:text;not null"`
	CreatedAt   time.Time `gorm:"created_at;type:timestamp;not null"`
	UpdatedAt   time.Time `gorm:"updated_at;type:timestamp;not null"`
}
