package handlers

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/helper"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/service"
	"net/http"
)

type NotificationHandler struct {
	notificationService service.NotificationService
}

func (n *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("id")
	if userId == "" {
		helper.BadRequestResponse(w, "id in query is required", nil)
		return
	}

	if err := n.notificationService.MarkAsRead(r.Context(), userId); err != nil {
		helper.InternalServerError(w, "Failed to mark notifications as read", err)
		return
	}

	notifications, err := n.notificationService.GetUserNotifications(r.Context(), userId)
	if err != nil {
		helper.InternalServerError(w, "Failed to retrieve notifications", err)
		return
	}

	helper.SuccessResponse(w, "Notifications marked as read", notifications)
}

func (n *NotificationHandler) GetUserNotifications(w http.ResponseWriter, r *http.Request) {
	userID := httprouter.ParamsFromContext(r.Context()).ByName("id")
	if userID == "" {
		helper.BadRequestResponse(w, "userid in path is required", nil)
		return
	}

	notifications, err := n.notificationService.GetUserNotifications(r.Context(), userID)
	if err != nil {
		helper.InternalServerError(w, "Failed to retrieve notifications", err)
		return
	}

	if len(notifications) == 0 {
		helper.SuccessResponse(w, "No notifications found", []any{})
		return
	}

	helper.SuccessResponse(w, "Notifications found", notifications)
}

func NewNotificationHandler(notificationService service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}
