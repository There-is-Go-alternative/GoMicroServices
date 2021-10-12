package usecase

import (
	"context"
	"fmt"

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
