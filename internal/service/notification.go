package service

import (
	"context"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/domain"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/repository"
)

type NotificationService interface {
	MarkAsRead(ctx context.Context, userId string) error
	GetUserNotifications(ctx context.Context, userId string) ([]*domain.Notification, error)
}

type notificationService struct {
	notificationRepository repository.NotificationRepository
}

func (n *notificationService) MarkAsRead(ctx context.Context, userId string) error {
	return n.notificationRepository.MarkAsRead(ctx, userId)
}
func (n *notificationService) GetUserNotifications(ctx context.Context, userId string) ([]*domain.Notification, error) {
	return n.notificationRepository.GetByUserId(ctx, userId)
}

func NewNotificationService(notificationRepository repository.NotificationRepository) NotificationService {
	return &notificationService{
		notificationRepository: notificationRepository,
	}
}
