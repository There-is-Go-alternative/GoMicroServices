package usecase

import (
	"context"

	"github.com/There-is-Go-alternative/GoMicroServices/chats/domain"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type database interface {
	GetAllChatsOfUser(ctx context.Context, user_id string) ([]*domain.Chat, error)
	GetChatByID(ctx context.Context, id domain.ChatID) (*domain.Chat, error)
	CreateChat(ctx context.Context, chats domain.Chat) error
	GetMessagesByChatID(ctx context.Context, id domain.ChatID) ([]*domain.Message, error)
	CreateMessage(ctx context.Context, message domain.Message) error
}

type UseCase struct {
	DB     database
	logger zerolog.Logger
}

func NewGetUseCase(db database) *UseCase {
	return &UseCase{
		DB:     db,
		logger: log.With().Str("service", "UseCase").Logger(),
	}
}
