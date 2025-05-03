package entity

import (
	"github.com/google/uuid"
	"time"
)

type OutboxMessage struct {
	Id        uuid.UUID `gorm:"id;type:uuid;default:gen_random_uuid();primary_key"`
	Message   string    `gorm:"message;type:jsonb;not null"`
	CreatedAt time.Time `gorm:"created_at;type:timestamp;not null"`
	UpdatedAt time.Time `gorm:"updated_at;type:timestamp;not null"`
	IsSent    bool      `gorm:"is_sent;type:boolean;not null"`
}

func (OutboxMessage) TableName() string {
	return "outbox_message"
}
