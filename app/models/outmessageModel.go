package models

import (
	pgorm "../../db/gorm"
)

type OutGoingMessage struct {
	Id           int64  `gorm:"primary_key" json:"id"`
	Type         int    `gorm:"column:type" json:"type"`
	Relay        string `gorm:"column:relay" json:"relay"`
	Timestamp    int64  `gorm:"column:timestamp" json:"timestamp"`
	Source       string `gorm:"column:source" json:"source"`
	SourceDevice int64  `gorm:"column:source_device" json:"source_device"`
	Message      []byte `gorm:"column:message" json:"message"`
	Content      []byte `gorm:"column:content" json:"content"`
	Destination  string `gorm:"column:destination" json:"destination"`
}

type OutGoingMessages struct {
	Messages []OutGoingMessage `gorm:"column:messages" json:"messages"`
	More     bool              `gorm:"column:more" json:"more"`
}

func (OutGoingMessage) TableName() string {
	return "public.messages"
}

func (m *OutGoingMessage) GetMessageForDevice(destination string, destinationDevice int64) (outGoingMessages OutGoingMessages) {
	messages := &outGoingMessages.Messages
	pgorm.MsgManager().Where("destination = ? and destination_device = ?", destination, destinationDevice).Find(messages)
	if len(*messages) > 100 {
		outGoingMessages.More = true
	}
	return
}

func (m *OutGoingMessage) SaveMessage() {
	pgorm.MsgManager().Save(m)
	return
}
