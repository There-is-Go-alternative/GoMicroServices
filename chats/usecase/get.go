package usecase

import (
	"context"

	"github.com/There-is-Go-alternative/GoMicroServices/chats/domain"
)

type GetAllChatsOfUserCmd func(ctx context.Context, user_id string) ([]*domain.Chat, error)

func (u UseCase) GetAllChatsOfUser() GetAllChatsOfUserCmd {
	return func(ctx context.Context, user_id string) ([]*domain.Chat, error) {
		u.logger.Info().Msg("Fetching all chats ...")
		defer u.logger.Info().Msg("All chats fetched !")
		//TODO: In the future, when the Database will be real, find another way to make this code return an empty array
		all_chats, err := u.DB.GetAllChatsOfUser(ctx, user_id)

		if all_chats == nil {
			return make([]*domain.Chat, 0), nil
		}
		return all_chats, err
	}
}

type GetChatByIdCmd func(ctx context.Context, id domain.ChatID) (*domain.Chat, error)

func (u UseCase) GetChatById() GetChatByIdCmd {
	return func(ctx context.Context, id domain.ChatID) (*domain.Chat, error) {
		u.logger.Info().Msgf("Fetching chat by id: %v", id)
		defer u.logger.Info().Msg("All chats fetched !")
		return u.DB.GetChatByID(ctx, id)
	}
}

type GetMessagesByChatIDCmd func(ctx context.Context, id domain.ChatID) ([]*domain.Message, error)

func (u UseCase) GetMessagesByChatID() GetMessagesByChatIDCmd {
	return func(ctx context.Context, id domain.ChatID) ([]*domain.Message, error) {
		u.logger.Info().Msgf("Fetching messages by chat id: %v", id)
		defer u.logger.Info().Msg("All chats fetched !")
		return u.DB.GetMessagesByChatID(ctx, id)
	}
}
