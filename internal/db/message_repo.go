package db

import (
	"gorm.io/gorm"
)

type MessageRepository interface {
	GetSentMessages() ([]Message, error)
	GetUnsentMessages(limit int) ([]Message, error)
	UpdateMessageStatus(messages []Message, status string) error
}

type GormMessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &GormMessageRepository{db: db}
}

func (r *GormMessageRepository) GetSentMessages() ([]Message, error) {
	var messages []Message
	if err := r.db.Where("status = ?", "sent").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *GormMessageRepository) GetUnsentMessages(limit int) ([]Message, error) {
	var messages []Message
	if err := r.db.Where("status != ?", "sent").Limit(limit).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *GormMessageRepository) UpdateMessageStatus(messages []Message, status string) error {
	for _, message := range messages {
		if err := r.db.Model(&message).Update("status", status).Error; err != nil {
			return err
		}
	}
	return nil
}
