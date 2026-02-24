package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/gateway/handlers"
	"net/http"
)

type NotificationRoute struct {
	notificationHandler *handlers.NotificationHandler
}

func (n *NotificationRoute) NotificationRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/v1/notifications/mark-read", n.notificationHandler.MarkAsRead)
	router.HandlerFunc(http.MethodGet, "/v1/notification/:id", n.notificationHandler.GetUserNotifications)
}

func NewNotificationRoute(notificationHandler *handlers.NotificationHandler) *NotificationRoute {
	return &NotificationRoute{
		notificationHandler: notificationHandler,
	}
}
