package service

import (
	"context"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/domain"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/dto"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/repository"
)

type MessageService interface {
	SendMessage(ctx context.Context, input *dto.MessageReq) (*dto.MessageResp, error)
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

func (m *messageService) toMessage(input *dto.MessageReq) *domain.Message {
	return &domain.Message{
		Content:  input.Content,
		Sender:   input.Sender,
		Receiver: input.Receiver,
	}
}

func (m *messageService) toMessageResp(input *domain.Message) *dto.MessageResp {
	return &dto.MessageResp{
		Id: input.Id,
	}
}

func NewMessageService(messageRepository repository.MessageRepository, unreadRepository repository.UnreadMessageRepository) MessageService {
	return &messageService{
		messageRepository: messageRepository,
		unreadRepository:  unreadRepository,
	}
}
