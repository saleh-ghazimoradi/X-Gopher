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
}

func NewMessageRoute(messageHandler *handlers.MessageHandler) *MessageRoute {
	return &MessageRoute{
		messageHandler: messageHandler,
	}
}
