package entity

import (
	"github.com/google/uuid"
	"time"
)

type OutboxMessage struct {
	Id        uuid.UUID `gorm:"id"`
	Message   string    `gorm:"message"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
	Status    string    `gorm:"status"`
}
