package service

import (
	"context"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/domain"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/dto"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/repository"
)

type MessageService interface {
	SendMessage(ctx context.Context, input *dto.MessageReq) (*dto.MessageResp, error)
	GetMessages(ctx context.Context, user1, user2 string, page int) ([]dto.MessageResp, error)
	GetUnreadSummary(ctx context.Context, receiverId string) (*dto.UnreadSummaryResp, error)
	MarkAsRead(ctx context.Context, receiverId, senderId string) error
}

type messageService struct {
	messageRepository repository.MessageRepository
	unreadRepository  repository.UnreadMessageRepository
}

func (m *messageService) SendMessage(ctx context.Context, input *dto.MessageReq) (*dto.MessageResp, error) {
	message := m.toMessage(input)
	if err := m.messageRepository.CreateMessage(ctx, message); err != nil {
		return nil, err
	}

	if err := m.unreadRepository.IncrementUnread(ctx, input.Sender, input.Receiver); err != nil {
		return nil, err
	}

	return m.toMessageResp(message), nil
}

func (m *messageService) GetMessages(ctx context.Context, user1, user2 string, page int) ([]dto.MessageResp, error) {

	const limit = 20
	skip := int64(page * limit)

	messages, err := m.messageRepository.
		GetMessagesBetween(ctx, user1, user2, skip, limit)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(messages)/2; i++ {
		j := len(messages) - 1 - i
		messages[i], messages[j] = messages[j], messages[i]
	}

	var resp []dto.MessageResp
	for _, msg := range messages {
		resp = append(resp, dto.MessageResp{Id: msg.Id, Content: msg.Content, Sender: msg.Sender, Receiver: msg.Receiver})
	}

	return resp, nil
}

func (m *messageService) GetUnreadSummary(ctx context.Context, receiverId string) (*dto.UnreadSummaryResp, error) {

	unreads, err := m.unreadRepository.GetUnreadByReceiver(ctx, receiverId)
	if err != nil {
		return nil, err
	}

	var total int
	var conversations []dto.UnreadConversation

	for _, u := range unreads {
		total += u.NumOfUnreadMessages

		conversations = append(conversations, dto.UnreadConversation{
			Id:                  u.Id,
			SenderId:            u.SenderId,
			ReceiverId:          u.ReceiverId,
			NumOfUnreadMessages: u.NumOfUnreadMessages,
			IsRead:              u.IsRead,
		})
	}

	return &dto.UnreadSummaryResp{
		Conversations: conversations,
		Total:         total,
	}, nil
}

func (m *messageService) MarkAsRead(ctx context.Context, receiverId, senderId string) error {
	return m.unreadRepository.MarkAsRead(ctx, receiverId, senderId)
}

func (m *messageService) toMessage(input *dto.MessageReq) *domain.Message {
	return &domain.Message{
		Content:  input.Content,
		Sender:   input.Sender,
		Receiver: input.Receiver,
	}
}

func (m *messageService) toMessageResp(input *domain.Message) *dto.MessageResp {
	return &dto.MessageResp{
		Id:       input.Id,
		Content:  input.Content,
		Sender:   input.Sender,
		Receiver: input.Receiver,
	}
}

func NewMessageService(messageRepository repository.MessageRepository, unreadRepository repository.UnreadMessageRepository) MessageService {
	return &messageService{
		messageRepository: messageRepository,
		unreadRepository:  unreadRepository,
	}
}
