package handlers

import (
	"github.com/saleh-ghazimoradi/X-Gopher/internal/dto"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/helper"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/service"
	"net/http"
	"strconv"
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

func (m *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	page, _ := strconv.Atoi(query.Get("page"))

	req := dto.GetMessagesQuery{
		Sender:   query.Get("sender"),
		Receiver: query.Get("receiver"),
		Page:     page,
	}

	v := helper.NewValidator()
	dto.ValidateGetMessagesQuery(v, &req)

	if !v.Valid() {
		helper.FailedValidationResponse(w, "Invalid query params")
		return
	}

	messages, err := m.messageService.GetMessages(r.Context(), req.Sender, req.Receiver, req.Page)
	if err != nil {
		helper.InternalServerError(w, "Failed to get messages", err)
		return
	}

	helper.SuccessResponse(w, "Messages retrieved", messages)
}

func (m *MessageHandler) GetUnreadSummary(w http.ResponseWriter, r *http.Request) {
	receiverId := r.URL.Query().Get("receiver")

	if receiverId == "" {
		helper.BadRequestResponse(w, "receiver is required", nil)
		return
	}

	resp, err := m.messageService.GetUnreadSummary(r.Context(), receiverId)
	if err != nil {
		helper.InternalServerError(w, "Failed to get unread messages", err)
		return
	}

	helper.SuccessResponse(w, "Unread summary", resp)
}

func (m *MessageHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	receiver := query.Get("receiver")
	sender := query.Get("sender")

	if receiver == "" || sender == "" {
		helper.BadRequestResponse(w, "receiver and sender are required", nil)
		return
	}

	if err := m.messageService.
		MarkAsRead(r.Context(), receiver, sender); err != nil {
		helper.InternalServerError(w, "Failed to mark as read", err)
		return
	}

	helper.SuccessResponse(w, "Messages marked as read", nil)
}

func NewMessageHandler(messageService service.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}
