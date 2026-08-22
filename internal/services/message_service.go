package services

import (
	"bugtracker/internal/dto"
	"bugtracker/internal/models"
	"bugtracker/internal/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageService struct {
	messageRepo *repositories.MessageRepo
}

func NewMessageService(messageRepo *repositories.MessageRepo) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
	}
}

func (c *MessageService) NewMessage(message dto.NewMessageDto) error {

	senderID, err := primitive.ObjectIDFromHex(message.SenderID)

	if err != nil {
		return err
	}

	cardID, err := primitive.ObjectIDFromHex(message.CardID)

	if err != nil {
		return err
	}

	if message.ReplyTo != "" {
		replyTo, err := primitive.ObjectIDFromHex(message.ReplyTo)

		if err != nil {
			return err
		}

		return c.messageRepo.NewMessage(&models.Message{
			SenderID: senderID,
			CardID:   cardID,
			ReplyTo:  replyTo,
			Content:  message.Content,
		})
	} else {
		return c.messageRepo.NewMessage(&models.Message{
			SenderID: senderID,
			CardID:   cardID,	
			Content:  message.Content,
		})
	}

}

func (c *MessageService) DeleteMessage(cardID string, userID string) {

}

func (c *MessageService) EditMessage() {

}
func (c *MessageService) GetMessages(cardID string) {

}
