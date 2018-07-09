package models

type IncomingMessage struct {
	Type                      int    `gorm:"column:type" json:"type"`
	Destination               string `gorm:"column:destination" json:"destination,omitempty"`
	DestinationDeviceId       int64  `gorm:"column:destinationDeviceId" json:"destinationDeviceId"`
	DestinationRegistrationId int    `gorm:"column:destinationRegistrationId" json:"destinationRegistrationId"`
	Body                      string `gorm:"column:body" json:"body"`
	Relay                     string `gorm:"column:relay" json:"relay,omitempty"`
	Silent                    bool   `gorm:"column:silent" json:"silent"`
	Content                   string `gorm:"column:content" json:"content"`
}

type IncomingMessageList struct {
	Messages    []IncomingMessage `gorm:"column:messages" json:"messages"`
	Destination string            `gorm:"column:destination" json:"destination,omitempty"`
	Relay       string            `gorm:"column:relay" json:"relay,omitempty"`
	Timestamp   int64             `gorm:"column:timestamp" json:"timestamp"`
}
