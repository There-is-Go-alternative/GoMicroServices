package database

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"github.com/There-is-Go-alternative/GoMicroServices/chats/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/chats/internal/xerrors"
	"github.com/pkg/errors"
)

var ChatsDefaultConf = &FirebaseConfig{
	CollectionName:    "chats",
	ServiceAdsKeyPath: "infra/database/FirebaseCredentials.json",
	BaseConfig: &firebase.Config{
		DatabaseURL: "https://gomicroservicedatabase-default-rtdb.firebaseio.com/",
	},
}

// Create adds a chat to the Firestore realtime database
func (m *FirebaseRealTimeDB) CreateChat(ctx context.Context, chat domain.Chat) error {

	err := m.DB.NewRef(fmt.Sprintf("%v/%v", m.Conf.CollectionName, chat.ID.String())).Set(ctx, chat)
	if err != nil {
		return err
	}
	return nil
}

// ByID Retrieve the info that match "id".
// Strict: As ID is the key of the map, return an error if not found
func (m *FirebaseRealTimeDB) ChatByID(ctx context.Context, ID domain.ChatID) (*domain.Chat, error) {
	var chat domain.Chat
	if err := m.DB.NewRef(fmt.Sprintf("%v/%v", m.Conf.CollectionName, ID)).Get(ctx, &chat); err != nil {
		return nil, err
	}

	if chat.ID.String() == "" {
		return nil, errors.Wrapf(
			xerrors.ErrorWithCode{Code: xerrors.ResourceNotFound, Err: xerrors.ChatNotFound}, "ID {%v}", ID,
		)
	}
	return &chat, nil
}

//Get all conversations of a user
func (m *FirebaseRealTimeDB) GetAllChatsOfUser(ctx context.Context, ID string) ([]*domain.Chat, error) {
	var chat map[string]*domain.Chat
	if err := m.DB.NewRef(m.Conf.CollectionName).OrderByKey().Get(ctx, &chat); err != nil {
		return nil, err
	}
	lst := make([]*domain.Chat, 0, len(chat))
	for _, a := range chat {
		for _, b := range a.UsersIDs {
			if b == ID {
				lst = append(lst, a)
				break
			}
		}
	}
	return lst, nil
}
