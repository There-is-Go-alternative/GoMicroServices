package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/There-is-Go-alternative/GoMicroServices/chats/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/chats/internal/xerrors"
)

type CreateChatCmd func(ctx context.Context, input CreateChatInput) (*domain.Chat, error)

type CreateChatInput struct {
	UsersIDs []string `json:"users_ids"`
}

func (u UseCase) CreateChat() CreateChatCmd {
	return func(ctx context.Context, input CreateChatInput) (*domain.Chat, error) {
		chat := domain.Chat{UsersIDs: input.UsersIDs}
		ChatID, err := domain.NewChatID()
		if err != nil {
			return nil, err
		}
		chat.ID = ChatID
		if !chat.Validate() {
			return nil, xerrors.ErrorWithCode{
				Code: xerrors.CodeInvalidData, Err: fmt.Errorf("invalid user chat data: %v", chat),
			}
		}
		err = u.DB.CreateChat(ctx, chat)
		return &chat, err
	}
}

type CreateMessageCmd func(ctx context.Context, input CreateMessageInput) (*domain.Message, error)

type CreateMessageInput struct {
	ChatID       string    `json:"chat_id"`
	Content      string    `json:"content"`
	SenderID     string    `json:"sender_id"`
	CreatedAt    time.Time `json:"created_at"`
	Attachements [][]byte  `json:"attachements"`
}

func (u UseCase) CreateMessage() CreateMessageCmd {
	return func(ctx context.Context, input CreateMessageInput) (*domain.Message, error) {
		message := domain.Message{ChatID: input.ChatID, Content: input.Content, SenderID: input.SenderID, CreatedAt: time.Now(), Attachements: input.Attachements}
		messageID, err := domain.NewMessageID()
		if err != nil {
			return nil, err
		}
		message.ID = messageID
		if !message.Validate() {
			return nil, xerrors.ErrorWithCode{
				Code: xerrors.CodeInvalidData, Err: fmt.Errorf("invalid message data: %v", message),
			}
		}
		err = u.DB.CreateMessage(ctx, message)
		return &message, err
	}
}
