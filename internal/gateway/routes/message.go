package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/gateway/handlers"
	"net/http"
)

type MessageRoute struct {
	messageHandler *handlers.MessageHandler
}

func (m *MessageRoute) MessageRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, "/v1/message/send", m.messageHandler.SendMessage)
	router.HandlerFunc(http.MethodGet, "/v1/messages", m.messageHandler.GetMessages)
	router.HandlerFunc(http.MethodGet, "/v1/messages/unread", m.messageHandler.GetUnreadSummary)
	router.HandlerFunc(http.MethodPatch, "/v1/messages/read", m.messageHandler.MarkAsRead)
}

func NewMessageRoute(messageHandler *handlers.MessageHandler) *MessageRoute {
	return &MessageRoute{
		messageHandler: messageHandler,
	}
}
