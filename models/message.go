package models

import (
	u "DummyMessengerAPI/utils"
	"github.com/jinzhu/gorm"
)

type Message struct {
	gorm.Model
	Body       string `json:"body"`
	SenderId   uint   `json:"sender_id"`
	ReceiverId uint   `json:"receiver_id"`
}

func (message *Message) Validate() (map[string]interface{}, bool) {
	if message.Body == "" {
		return u.Message(false, "Empty message body"), false
	}
	if message.ReceiverId <= 0 {
		return u.Message(false, "User is not recognized - ReceiverId <= 0"), false
	}
	if message.SenderId <= 0 {
		return u.Message(false, "User is not recognized - SenderId <= 0"), false
	}
	if message.SenderId == message.ReceiverId {
		return u.Message(false, "Attempt to send message to yourself"), false
	}
	return u.Message(true, "Requirement passed"), true
}

func (message *Message) Create() map[string]interface{} {
	if resp, ok := message.Validate(); !ok {
		return resp
	}
	GetDB().Create(message)
	return u.Message(true, "Message has been created")
}

func GetLastMessage(senderId, receiverId uint) *Message {
	lastReceivedMessage := &Message{}
	lastSendedMessage := &Message{}
	err1 := GetDB().Table("messages").Where("sender_id = ? and receiver_id = ?", senderId, receiverId).Last(lastSendedMessage).Error
	err2 := GetDB().Table("messages").Where("sender_id = ? and receiver_id = ?", receiverId, senderId).Last(lastReceivedMessage).Error
	if err1 == gorm.ErrRecordNotFound && err2 == err1 {
		return nil
	}
	if lastReceivedMessage.CreatedAt.Unix() < lastSendedMessage.CreatedAt.Unix() {
		return lastSendedMessage
	}
	return lastReceivedMessage
}

func GetMessages(senderId, receiverId uint) []*Message {
	sendedMessages := make([]*Message, 0)
	receivedMessages := make([]*Message, 0)
	err1 := GetDB().Table("messages").Where("sender_id = ? and receiver_id = ?", senderId, receiverId).Find(&sendedMessages).Error
	err2 := GetDB().Table("messages").Where("sender_id = ? and receiver_id = ?", receiverId, senderId).Find(&receivedMessages).Error
	if err1 == gorm.ErrRecordNotFound && err1 == err2 {
		return nil
	}

	n := len(sendedMessages) + len(receivedMessages)
	messagesBody := make([]*Message, n)
	for i, p1, p2 := 0, 0, 0; i < n; i++ {
		if p1 == len(sendedMessages) {
			messagesBody[i] = receivedMessages[p2]
			p2++
		} else if p2 == len(receivedMessages) {
			messagesBody[i] = sendedMessages[p1]
			p1++
		} else if sendedMessages[p1].CreatedAt.Unix() > receivedMessages[p2].CreatedAt.Unix() {
			messagesBody[i] = receivedMessages[p2]
			p2++
		} else {
			messagesBody[i] = sendedMessages[p1]
			p1++
		}
	}
	return messagesBody
}
