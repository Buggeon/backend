package services

type NotificationService struct {
}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (n *NotificationService) SendNotificationByEmail(email string, message string) {

}

func (n *NotificationService) SendNotificationByTelegram(email string, message string) {

}
