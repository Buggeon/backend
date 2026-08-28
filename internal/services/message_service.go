package services

import (
	"bugtracker/internal/dto"
	"bugtracker/internal/models"
	"bugtracker/internal/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageService struct {
	messageRepo *repositories.MessageRepo
	cardRepo    *repositories.CardRepo
}

func NewMessageService(messageRepo *repositories.MessageRepo, cardRepo *repositories.CardRepo) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
		cardRepo:    cardRepo,
	}
}

func (c *MessageService) NewMessage(message dto.NewMessageDto) error {

	senderID, err := primitive.ObjectIDFromHex(message.SenderID)

	if err != nil {
		return err
	}

	cardID, err := primitive.ObjectIDFromHex(message.CardID)
	var messageID primitive.ObjectID

	if err != nil {
		return err
	}

	if message.ReplyTo != "" {
		replyTo, err := primitive.ObjectIDFromHex(message.ReplyTo)

		if err != nil {
			return err
		}

		id, err := c.messageRepo.NewMessage(&models.Message{
			SenderID: senderID,
			CardID:   cardID,
			ReplyTo:  replyTo,
			Content:  message.Content,
		})

		if err != nil {
			return err
		}

		messageID = id

	} else {
		id, err := c.messageRepo.NewMessage(&models.Message{
			SenderID: senderID,
			CardID:   cardID,
			Content:  message.Content,
		})

		if err != nil {
			return err
		}

		messageID = id
	}

	return c.cardRepo.AddMessage(cardID, messageID)

}

func (c *MessageService) DeleteMessage(cardID string, userID string) {

}

func (c *MessageService) EditMessage() {

}

func (c *MessageService) GetMessages(cardID string) ([]models.Message, error) {

	objID, err := primitive.ObjectIDFromHex(cardID)

	if err != nil {
		return nil, err
	}

	return c.messageRepo.GetMessages(objID)

}

func (c *MessageService) GetMessage(messageID string) (models.Message, error) {

	messageObjID, err := primitive.ObjectIDFromHex(messageID)

	if err != nil {
		return models.Message{}, nil
	}

	return c.messageRepo.GetMessage(messageObjID)

}
