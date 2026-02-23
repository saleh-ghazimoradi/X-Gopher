package handlers

import (
	"github.com/saleh-ghazimoradi/X-Gopher/internal/dto"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/helper"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/service"
	"net/http"
)

type MessageHandler struct {
	messageService service.MessageService
}

func (m *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var payload dto.MessageReq
	if err := helper.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, "Invalid given payload", err)
		return
	}

	v := helper.NewValidator()
	dto.ValidateMessageReq(v, &payload)
	if !v.Valid() {
		helper.FailedValidationResponse(w, "Invalid given payload")
		return
	}

	message, err := m.messageService.SendMessage(r.Context(), &payload)
	if err != nil {
		helper.InternalServerError(w, "Failed to create message", err)
		return
	}

	helper.CreatedResponse(w, "Message successfully sent", message.Id)
}

func NewMessageHandler(messageService service.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}
