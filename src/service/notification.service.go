package service

import (
	"github.com/peam1146/mcv-notifier/src/model"
	"gorm.io/gorm"
)

type NotificationService interface {
	SaveNotifications(notifications []*model.Notification) ([]*model.Notification, error)
	SaveNotification(notification *model.Notification) (*model.Notification, error)
}

type notificationService struct {
	db *gorm.DB
}

func NewNotificationService(db *gorm.DB) NotificationService {
	return &notificationService{db}
}

func (ns *notificationService) SaveNotification(notification *model.Notification) (*model.Notification, error) {
	if err := ns.db.Create(&notification).Error; err != nil {
		return nil, err
	}
	return notification, nil
}

func (ns *notificationService) SaveNotifications(notifications []*model.Notification) ([]*model.Notification, error) {
	var latestNotifications []*model.Notification
	for _, notification := range notifications {
		n, err := ns.SaveNotification(notification)
		if err != nil && err.Error() == "UNIQUE constraint failed: notifications.id" {
			continue
		}

		if err != nil {
			return nil, err
		}
		latestNotifications = append(latestNotifications, n)
	}
	return latestNotifications, nil
}
